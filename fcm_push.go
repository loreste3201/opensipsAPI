package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var fcmMessageJSON map[string]interface{}
var fcmKeys map[string]map[string]string

var fcmKey = ""
var fcmURL = ""
var fcmAvailable = false
var httpClient = &http.Client{}

//InitializeFCM ()
func InitializeFCM() {

	params := make([]interface{}, 0)
	ResultSet, ok := GetAllRows("SELECT * FROM `fcm_keys` WHERE `enabled`=1", params, "default")
	if !ok {
		Log("No fcm key in DB")
		return
	}
	Log(fmt.Sprintf("%d fcm keys loaded", len(ResultSet)))
	fcmKeys = make(map[string]map[string]string)
	for l := 0; l < len(ResultSet); l++ {
		fcmKeys[ResultSet[l]["name"]] = ResultSet[l]
	}
	Debug(fmt.Sprintf("fcmKeys::%v", fcmKeys))

	/*
		value, ok := Config["FCM_KEY"]
		if ok {
			fcmKey = value
		} else {
			Error("Missing configuration FCM_KEY")
			return
		}
		value, ok = Config["FCM_URL"]
		if ok {
			fcmURL = value
		} else {
			Error("Missing configuration FCM_URL")
			return
		}
	*/
	fcmMessageJSON = make(map[string]interface{})
	fcmAvailable = true
	return
}

//SendFCMPush (UserID as string, FCM-DeviceToken as string, Data as map[string]interface{} will be send as Data value in push JSON, HighPriority bool)
func SendFCMPush(UserID string, DeviceToken string, Title string, Data map[string]interface{}, HighPriority bool, SlientPush bool, DeviceType string, AppVersion string, FcmVersion string) bool {

	if fcmKeys[FcmVersion]["fcm_key"] == "" {
		FcmVersion = "default"
	}
	Log(fmt.Sprintf("Sending Push to Fcm Version::%v", FcmVersion))
	fcmURL = fcmKeys[FcmVersion]["fcm_url"]
	fcmKey = fcmKeys[FcmVersion]["fcm_key"]

	logText := fmt.Sprintf("::\tTo:%s\tTitle:%s", DeviceToken, Title)
	fcmMessageJSON = make(map[string]interface{})
	Debug(logText)
	Debug(fmt.Sprintf("SilentPush::%v", SlientPush))
	params := make([]interface{}, 9)
	params[0] = "FCM"
	params[1] = UserID
	params[2] = DeviceToken
	params[4] = "FAIL"
	if !fcmAvailable {
		params[6] = "FCM Service not available"
		Error("SendFCMPush:: FCM Service not available")
		return false
	}
	fcmMessageJSON["to"] = DeviceToken
	if len(Data) > 0 {
		fcmMessageJSON["content_available"] = true
		fcmMessageJSON["data"] = Data
	} else {
		fcmMessageJSON["content_available"] = false
	}
	if HighPriority {
		fcmMessageJSON["priority"] = "high"
	} else {
		fcmMessageJSON["priority"] = "medium"
	}
	if SlientPush == false {
		messageNotification := make(map[string]interface{})
		if strings.ToLower(DeviceType) == "ios" {

			if Data["notification_type"] == "offline_message" || Data["notification_type"] == "friend_request_received" || Data["notification_type"] == "friend_request_rejected" || Data["notification_type"] == "friend_request_cancelled" || Data["notification_type"] == "friend_request_accepted" {
				messageNotification["title"] = Title
				messageNotification["body"] = Data["message"]
			} else {
				s := strings.Split(AppVersion, " ")
				s1 := strings.Split(s[0], ".")

				leftVersion := ""
				for j := 0; j < len(s1); j++ {
					leftVersion = leftVersion + s1[j]
				}
				leftVersionInt, _ := strconv.Atoi(leftVersion)
				if leftVersionInt <= 101 {
					messageNotification["title"] = Title
				} else {
					messageNotification["body"] = Title
				}
			}

		} else {
			messageNotification["title"] = Title
		}
		messageNotification["sound"] = "default"
		//messageNotification["badge"] = 0
		fcmMessageJSON["notification"] = messageNotification
	}

	jsonString, err := json.Marshal(fcmMessageJSON)
	if err != nil {
		Error(logText + "\tcan't Marshal JSON, Error:" + err.Error())
		params[6] = "can't Marshal JSON, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	params[3] = jsonString
	httpRequest, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonString))
	defer httpRequest.Body.Close()
	if err != nil {
		Error(logText + "\tcan't create HTTP Request, Error:" + err.Error())
		params[6] = "can't create HTTP Request, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	httpRequest.Header.Set("Authorization", "key="+fcmKey)
	httpRequest.Header.Set("Content-Type", "application/json")
	//make HTTP Request
	//httpClient := &http.Client{}
	Response, err := httpClient.Do(httpRequest)
	if err != nil {
		Error(logText + "\tcan't make HTTP Request, Error:" + err.Error())
		params[6] = "can't make HTTP Request, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	defer Response.Body.Close()

	Debug(logText + "\tFCM Response Status:" + string(Response.StatusCode) + "\t" + Response.Status)
	Debug(fmt.Sprintf(logText, "\tFCM Response Headers:%v", Response.Header))
	params[5] = Response.StatusCode
	params[6] = Response.Status
	params[7] = fmt.Sprintf("%v", Response.Header)
	if Response.StatusCode != 200 {
		addFCMCDR(params)
		return false
	}
	ResponseBody, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		Error(logText + "\tcan't read Response Body, Error:" + err.Error())
		params[8] = "can't read Response Body, Error:" + err.Error()
		return false
	}
	params[8] = ResponseBody
	Debug(logText + "\tFCM Response Body:" + string(ResponseBody))
	if strings.Contains(strings.ToLower(Response.Header.Get("Content-Type")), "application/json") {
		ResponseBodyJSON := make(map[string]interface{})
		if err := json.Unmarshal(ResponseBody, &ResponseBodyJSON); err != nil {
			Error(logText + "\tcan't convert Response Body to JSON, Error:" + err.Error())
			//params[9] = "can't convert Response Body to JSON, Error:" + err.Error()
			addFCMCDR(params)
			return false
		}
		Debug("Response: " + fmt.Sprintf("%v", ResponseBodyJSON))
		if InterfaceToInt(ResponseBodyJSON["success"]) > 0 {
			params[4] = "SUCCESS"
		}
	}

	addFCMCDR(params)
	return true
}

//SendFCMPushForCall (UserID as string, FCM-DeviceToken as string, Data as map[string]interface{} will be send as Data value in push JSON, HighPriority bool)
func SendFCMPushForCall(UserID string, DeviceToken string, Title string, Data map[string]interface{}, HighPriority bool, SlientPush bool, DeviceType string, AppVersion string, FcmVersion string) bool {

	if fcmKeys[FcmVersion]["fcm_key"] == "" {
		FcmVersion = "default"
	}
	Log(fmt.Sprintf("Sending Push to Fcm Version::%v", FcmVersion))
	fcmURL = fcmKeys[FcmVersion]["fcm_url"]
	fcmKey = fcmKeys[FcmVersion]["fcm_key"]

	logText := fmt.Sprintf("SendFCMPush::\tTo:%s\tTitle:%s", DeviceToken, Title)
	apsJSON := make(map[string]interface{})
	fcmMessageJSON = make(map[string]interface{})
	Debug(logText)
	Debug(fmt.Sprintf("SlientPush::%v", SlientPush))
	params := make([]interface{}, 9)
	params[0] = "FCM"
	params[1] = UserID
	params[2] = DeviceToken
	params[4] = "FAIL"
	if !fcmAvailable {
		params[6] = "FCM Service not available"
		Error("SendFCMPush:: FCM Service not available")
		return false
	}
	fcmMessageJSON["to"] = DeviceToken
	if len(Data) > 0 {
		//fcmMessageJSON["content_available"] = true

		alertJSON := make(map[string]interface{})
		alertJSON["body"] = Title
		apsJSON["custom"] = Data
		apsJSON["alert"] = alertJSON
		apsJSON["content_available"] = true
		apsJSON["notification_type"] = Data["notification_type"]

		fcmMessageJSON["data"] = apsJSON
	} else {
		fcmMessageJSON["content_available"] = false
	}
	if HighPriority {
		fcmMessageJSON["priority"] = "high"
	} else {
		fcmMessageJSON["priority"] = "medium"
	}
	if SlientPush == false {
		messageNotification := make(map[string]interface{})
		if strings.ToLower(DeviceType) == "ios" {

			if Data["notification_type"] == "offline_message" || Data["notification_type"] == "friend_request_received" || Data["notification_type"] == "friend_request_rejected" || Data["notification_type"] == "friend_request_cancelled" {
				messageNotification["title"] = Title
				messageNotification["body"] = Data["message"]
			} else {
				s := strings.Split(AppVersion, " ")
				s1 := strings.Split(s[0], ".")

				leftVersion := ""
				for j := 0; j < len(s1); j++ {
					leftVersion = leftVersion + s1[j]
				}
				leftVersionInt, _ := strconv.Atoi(leftVersion)
				if leftVersionInt <= 101 {
					messageNotification["title"] = Title
				} else {
					messageNotification["body"] = Title
				}
			}

		} else {
			messageNotification["title"] = Title
		}
		messageNotification["sound"] = "default"
		//messageNotification["badge"] = 0
		fcmMessageJSON["notification"] = messageNotification
	}

	jsonString, err := json.Marshal(fcmMessageJSON)
	if err != nil {
		Error(logText + "\tcan't Marshal JSON, Error:" + err.Error())
		params[6] = "can't Marshal JSON, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	params[3] = jsonString
	httpRequest, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonString))
	defer httpRequest.Body.Close()
	if err != nil {
		Error(logText + "\tcan't create HTTP Request, Error:" + err.Error())
		params[6] = "can't create HTTP Request, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	httpRequest.Header.Set("Authorization", "key="+fcmKey)
	httpRequest.Header.Set("Content-Type", "application/json")
	//make HTTP Request
	//httpClient := &http.Client{}
	Response, err := httpClient.Do(httpRequest)
	if err != nil {
		Error(logText + "\tcan't make HTTP Request, Error:" + err.Error())
		params[6] = "can't make HTTP Request, Error:" + err.Error()
		addFCMCDR(params)
		return false
	}
	defer Response.Body.Close()

	Debug(logText + "\tFCM Response Status:" + string(Response.StatusCode) + "\t" + Response.Status)
	Debug(fmt.Sprintf(logText, "\tFCM Response Headers:%v", Response.Header))
	params[5] = Response.StatusCode
	params[6] = Response.Status
	params[7] = fmt.Sprintf("%v", Response.Header)
	if Response.StatusCode != 200 {
		addFCMCDR(params)
		return false
	}
	ResponseBody, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		Error(logText + "\tcan't read Response Body, Error:" + err.Error())
		params[8] = "can't read Response Body, Error:" + err.Error()
		return false
	}
	params[8] = ResponseBody
	Debug(logText + "\tFCM Response Body:" + string(ResponseBody))
	if strings.Contains(strings.ToLower(Response.Header.Get("Content-Type")), "application/json") {
		ResponseBodyJSON := make(map[string]interface{})
		if err := json.Unmarshal(ResponseBody, &ResponseBodyJSON); err != nil {
			Error(logText + "\tcan't convert Response Body to JSON, Error:" + err.Error())
			//params[9] = "can't convert Response Body to JSON, Error:" + err.Error()
			addFCMCDR(params)
			return false
		}
		Debug("Response: " + fmt.Sprintf("%v", ResponseBodyJSON))
		if InterfaceToInt(ResponseBodyJSON["success"]) > 0 {
			params[4] = "SUCCESS"
		}
	}

	addFCMCDR(params)
	return true
}

func addFCMCDR(params []interface{}) {
	UpdateDB("INSERT INTO `push_sent`(`push_type`,`user_id`,`device_token`,`sent_json`,`status`,`status_code`,`status_description`,`recieved_headers`,`recieved_body`,`created_datetime`)VALUES(?,?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())", params, "default")
}
