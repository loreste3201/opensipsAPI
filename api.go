package main

import (
	"fmt"
	"reflect"
	"strings"
)

//var dynamicAPI map[string]interface{}
var dynamicAPI map[string]map[string]string

//InitializDynamicAPIs Load dynamic APIs from DB
func InitializDynamicAPIs() {
	dynamicAPI = make(map[string]map[string]string)
	params := make([]interface{}, 0)
	dynamicAPIs, ok := GetAllRows("SELECT * FROM `dynamic_apis`", params, "default")
	if !ok {
		Warning("can't load dynamic APIs from DB")
		return
	}
	for i := 0; i < len(dynamicAPIs); i++ {
		apiData := dynamicAPIs[i]
		apiEndpoint := apiData["endpoint_name"]
		Debug("Loading API " + apiEndpoint)
		/*
			thisAPI := make(map[string]string)
			thisAPI["api_name"] = apiData["api_name"]
			thisAPI["description"] = apiData["description"]
			thisAPI["required_params"] = apiData["required_params"]
			thisAPI["sql"] = apiData["sql"]
			thisAPI["sql_params"] = apiData["sql_params"]
			thisAPI["api_type"] = apiData["api_type"]
			dynamicAPI[apiEndpoint] = thisAPI
		*/
		dynamicAPI[apiEndpoint] = apiData
	}
	Info(fmt.Sprintf("%d dynamic APIs loaded", len(dynamicAPIs)))
	return
}

//API (ResponseJSON map[string]interface{}, RequestJSON map[string]interface{}) Handle all APIs
func API(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {
	endpoint := InterfaceToString(RequestJSON["endpoint"])
	Debug("API: Endpoint : " + endpoint)
	var apiData map[string]string
	apiData, ok := dynamicAPI[endpoint]
	if ok {
		switch  RequestJSON["request_method"]{
		case "GET":
			apiGet(RequestJSON, ResponseJSON, apiData)
			break
		case "PUT":
			apiPUT(RequestJSON, ResponseJSON, apiData)
			break
		case "UPDATE":
			apiUpdate(RequestJSON, ResponseJSON, apiData)
			break
		default:
			Error("Invalid Dynamic API Type for endpoint " + endpoint)
			ResponseJSON["status"] = 500
			ResponseJSON["message"] = "Server configuration error"
			break
		}
	} else {
		if RequestJSON["API_TYPE"] == "api" {
			switch InterfaceToString(RequestJSON["api_version"]) {
			case "v1":
				CallStaticAPIsv1(RequestJSON, ResponseJSON, AuthData)
				break
			default:
				CallStaticAPIs(RequestJSON, ResponseJSON, AuthData)
				break
			}
		} else {
			AdminAPI(RequestJSON, ResponseJSON, AuthData)
		}

	}
}

//AppAPIs (ResponseJSON map[string]interface{}, RequestJSON map[string]interface{}) Handle all application APIs
func AppAPIs(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}) {

}

func apiGet(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, APIData map[string]string) {
	Query, ok := APIData["sql"]
	if !ok {
		ResponseJSON["status"] = 500
		ResponseJSON["message"] = "Server configuration error"
		Error("appGet: Invalid API data for " + RequestJSON["endpoint"].(string) + "\t sql not found")
		return
	}

	params := make([]interface{}, 0)
	requiredParams, ok := APIData["sql_params"]
	if (ok) && (len(requiredParams) > 2) {
		RequestData, ok := RequestJSON["data"]
		if !ok {
			ResponseJSON["status"] = 486
			ResponseJSON["message"] = "Missing data"
			return
		}
		if reflect.ValueOf(RequestData).Kind() != reflect.Map {
			ResponseJSON["status"] = 486
			ResponseJSON["message"] = "inavlid Data, should be a valid JSON Object"
			return
		}
		DataMap := RequestData.(map[string]interface{})

		paramsArray := strings.Split(requiredParams, ",")
		for i := 0; i < len(paramsArray); i++ {
			paramName := paramsArray[i]
			if paramName != "" {
				value, ok := DataMap[paramName]
				if !ok {
					ResponseJSON["status"] = 400
					ResponseJSON["message"] = "Missing parameter " + paramName
					return
				}
				params = append(params, value)
			}
		}
	}
	conditionState := "WHERE"
	if strings.Index(strings.ToUpper(Query), conditionState) > -1 {
		conditionState = "AND"
	}
	strSQL, ok := AdvancedSQL(RequestJSON, ResponseJSON, conditionState)
	if !ok {
		return
	}
	Query += strSQL
	rowsData, haveRows := GetAllRows(Query, params, "default")
	if !haveRows {
		ResponseJSON["status"] = 404
		ResponseJSON["message"] = "No Data"
	}

	returnType, ok := APIData["return_type"]
	if !ok {
		returnType = "SINGLE_ROW"
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = "Success"
	if returnType == "singleRow" {
		if haveRows {
			ResponseJSON["data"] = rowsData[0]
		} else {
			ResponseJSON["data"] = ""
		}
	} else {
		if haveRows {
			ResponseJSON["data"] = rowsData
		} else {
			ResponseJSON["data"] = make([]string, 0)
		}
	}
	return
}

func apiPUT(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, APIData map[string]string) {
	Query, ok := APIData["sql"]
	if !ok {
		ResponseJSON["status"] = 500
		ResponseJSON["message"] = "Server configuration error"
		Error("appGet: Invalid API data for " + RequestJSON["endpoint"].(string) + "\t sql not found")
		return
	}
	params := make([]interface{}, 0)
	requiredParams, ok := APIData["sql_params"]
	if ok {
		paramsArray := strings.Split(requiredParams, ",")
		for i := 0; i < len(paramsArray); i++ {
			paramName := paramsArray[i]
			value, ok := RequestJSON[paramName]
			if !ok {
				ResponseJSON["status"] = 400
				ResponseJSON["message"] = "Missing parameter " + paramName
				return
			}
			params = append(params, value)
		}
	}
	lastInsertID, ok := UpdateDB(Query, params, "default")
	if ok {
		ResponseJSON["status"] = 200
		ResponseJSON["message"] = "Success"
		ResponseJSON["data"] = fmt.Sprintf("Record Added ID %d", lastInsertID)
		return
	}
	ResponseJSON["status"] = 500
	ResponseJSON["message"] = "can't add data"
}

func apiUpdate(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, APIData map[string]string) {
	Query, ok := APIData["sql"]
	if !ok {
		ResponseJSON["status"] = 500
		ResponseJSON["message"] = "Server configuration error"
		Error("appGet: Invalid API data for " + RequestJSON["endpoint"].(string) + "\t sql not found")
		return
	}
	params := make([]interface{}, 0)
	requiredParams, ok := APIData["sql_params"]
	if ok {
		paramsArray := strings.Split(requiredParams, ",")
		for i := 0; i < len(paramsArray); i++ {
			paramName := paramsArray[i]
			value, ok := RequestJSON[paramName]
			if !ok {
				ResponseJSON["status"] = 400
				ResponseJSON["message"] = "Missing parameter " + paramName
				return
			}
			params = append(params, value)
		}
	}
	lastInsertID, ok := UpdateDB(Query, params, "default")
	if ok {
		ResponseJSON["status"] = 200
		ResponseJSON["message"] = "Success"
		ResponseJSON["data"] = fmt.Sprintf("%d Records Updated", lastInsertID)
		return
	}
	ResponseJSON["status"] = 400
	ResponseJSON["message"] = "can't update"
}

