package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strconv"
	"time"

	"github.com/oschwald/geoip2-golang"
)

// CallStaticAPIs ReInitialize Static APIs
func CallStaticAPIs(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {
	switch RequestJSON["endpoint"] {
	case "Login" + InterfaceToString(RequestJSON["api_version"]):
		LoginAPI(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetProfileInformation" + InterfaceToString(RequestJSON["api_version"]):
		GetProfileInformation(RequestJSON, ResponseJSON, AuthData)
		break
	case "ListOpSockets" + InterfaceToString(RequestJSON["api_version"]):
		ListOpSockets(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddAdmin" + InterfaceToString(RequestJSON["api_version"]):
		AddAdmin(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateAdmin" + InterfaceToString(RequestJSON["api_version"]):
		UpdateAdmin(RequestJSON, ResponseJSON, AuthData)
		break
	case "DeleteAdmin" + InterfaceToString(RequestJSON["api_version"]):
		DeleteAdmin(RequestJSON, ResponseJSON, AuthData)
		break
	case "ListAdmin" + InterfaceToString(RequestJSON["api_version"]):
		ListAdmin(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddDataCenter" + InterfaceToString(RequestJSON["api_version"]):
		AddDispatcherGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateDataCenter" + InterfaceToString(RequestJSON["api_version"]):
		UpdateDispatcherGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveDataCenter" + InterfaceToString(RequestJSON["api_version"]):
		DeleteDispatcherGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetDataCenters" + InterfaceToString(RequestJSON["api_version"]):
		ListDispatcherGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetDataCenter" + InterfaceToString(RequestJSON["api_version"]):
		ListDispatcherGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddGateway" + InterfaceToString(RequestJSON["api_version"]):
		AddDispatcher(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateGateway" + InterfaceToString(RequestJSON["api_version"]):
		UpdateDispatcher(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveGateway" + InterfaceToString(RequestJSON["api_version"]):
		DeleteDispatcher(RequestJSON, ResponseJSON, AuthData)
		break
	case "DisableGateway" + InterfaceToString(RequestJSON["api_version"]):
		DisableGateway(RequestJSON, ResponseJSON, AuthData)
		break
	case "EnableGateway" + InterfaceToString(RequestJSON["api_version"]):
		EnableGateway(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetGateways" + InterfaceToString(RequestJSON["api_version"]):
		ListDispatcher(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetGateway" + InterfaceToString(RequestJSON["api_version"]):
		ListDispatcher(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddRoute" + InterfaceToString(RequestJSON["api_version"]):
		AddSubscriberGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateRoute" + InterfaceToString(RequestJSON["api_version"]):
		UpdateSubscriberGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveRoute" + InterfaceToString(RequestJSON["api_version"]):
		RemoveRoute(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetRoutes" + InterfaceToString(RequestJSON["api_version"]):
		ListSubscriberGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetRoutesByLocation" + InterfaceToString(RequestJSON["api_version"]):
		ListSubscriberGroup(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddSubscriber" + InterfaceToString(RequestJSON["api_version"]):
		AddSubscriber(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateSubscriber" + InterfaceToString(RequestJSON["api_version"]):
		UpdateSubscriber(RequestJSON, ResponseJSON, AuthData)
		break
	case "DeleteSubscriber" + InterfaceToString(RequestJSON["api_version"]):
		DeleteSubscriber(RequestJSON, ResponseJSON, AuthData)
		break
	case "ListSubscriber" + InterfaceToString(RequestJSON["api_version"]):
		ListSubscriber(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocationIps" + InterfaceToString(RequestJSON["api_version"]):
		ListPermission(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocationIp" + InterfaceToString(RequestJSON["api_version"]):
		ListPermission(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddLocationIp" + InterfaceToString(RequestJSON["api_version"]):
		AddPermission(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateLocationIp" + InterfaceToString(RequestJSON["api_version"]):
		UpdatePermission(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveLocationIp" + InterfaceToString(RequestJSON["api_version"]):
		DeletePermission(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocationAnis" + InterfaceToString(RequestJSON["api_version"]):
		GetLocationAnis(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocationAni" + InterfaceToString(RequestJSON["api_version"]):
		GetLocationAnis(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddLocationAni" + InterfaceToString(RequestJSON["api_version"]):
		AddLocationAni(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateLocationAni" + InterfaceToString(RequestJSON["api_version"]):
		UpdateLocationAni(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveLocationAni" + InterfaceToString(RequestJSON["api_version"]):
		RemoveLocationAni(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddCarrier" + InterfaceToString(RequestJSON["api_version"]):
		AddCarrier(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateCarrier" + InterfaceToString(RequestJSON["api_version"]):
		UpdateCarrier(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveCarrier" + InterfaceToString(RequestJSON["api_version"]):
		RemoveCarrier(RequestJSON, ResponseJSON, AuthData)
		break
	case "DisableCarrier" + InterfaceToString(RequestJSON["api_version"]):
		DisableCarrier(RequestJSON, ResponseJSON, AuthData)
		break
	case "EnableCarrier" + InterfaceToString(RequestJSON["api_version"]):
		EnableCarrier(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetCarriers" + InterfaceToString(RequestJSON["api_version"]):
		GetCarriers(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetCarrier" + InterfaceToString(RequestJSON["api_version"]):
		GetCarriers(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddCarrierIp" + InterfaceToString(RequestJSON["api_version"]):
		AddCarrierIp(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateCarrierIp" + InterfaceToString(RequestJSON["api_version"]):
		UpdateCarrierIp(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveCarrierIp" + InterfaceToString(RequestJSON["api_version"]):
		RemoveCarrierIp(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetCarrierIps" + InterfaceToString(RequestJSON["api_version"]):
		GetCarrierIps(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetCarrierIp" + InterfaceToString(RequestJSON["api_version"]):
		GetCarrierIps(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocations" + InterfaceToString(RequestJSON["api_version"]):
		GetLocations(RequestJSON, ResponseJSON, AuthData)
		break
	case "GetLocation" + InterfaceToString(RequestJSON["api_version"]):
		GetLocations(RequestJSON, ResponseJSON, AuthData)
		break
	case "AddLocation" + InterfaceToString(RequestJSON["api_version"]):
		AddLocation(RequestJSON, ResponseJSON, AuthData)
		break
	case "UpdateLocation" + InterfaceToString(RequestJSON["api_version"]):
		UpdateLocation(RequestJSON, ResponseJSON, AuthData)
		break
	case "RemoveLocation" + InterfaceToString(RequestJSON["api_version"]):
		RemoveLocation(RequestJSON, ResponseJSON, AuthData)
		break
	default:
		ResponseJSON["status"] = 404
		ResponseJSON["message"] = "UnKnown API Call"
		break
	}
}

// LoginAPI for users
func LoginAPI(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	if "POST" != RequestJSON["request_method"] {
		ResponseJSON["status"] = 404
		ResponseJSON["message"] = GetMessageByID("en", 238) //"Invalid request method"
		return
	}

	ok, _, userID := GetUserIdentification(RequestJSON, ResponseJSON)
	if !ok {
		return
	}

	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 6) //"password required"
		return
	}

	params := make([]interface{}, 2)
	params[0] = userID
	params[1] = userID
	UserData, ok := GetSingleRow("SELECT * FROM `users` WHERE (`email` = ? OR `username` = ?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 401
		ResponseJSON["message"] = GetMessageByID("en", 8) //"Invalid Credentials"
		return
	}

	PasswordBytes := md5.Sum([]byte(InterfaceToString(password)))
	Password := fmt.Sprintf("%x", PasswordBytes)

	if UserData["password"] != Password {
		ResponseJSON["status"] = 401
		ResponseJSON["message"] = GetMessageByID("en", 8) //"Invalid Credentials"
		return
	}

	if UserData["active"] == "NO" {
		ResponseJSON["status"] = 403
		ResponseJSON["message"] = GetMessageByID("en", 185) //"User Verification required"
		return
	}

	NenoSeconds := time.Now().UnixNano()
	TokenString := UserData["email"] + strconv.FormatInt(NenoSeconds, 10)
	Token := md5.Sum([]byte(TokenString))
	AuthToken := fmt.Sprintf("%x", Token)
	Log("AuthToken " + AuthToken)

	db, err := geoip2.Open(Config["Country_Codes_File"])
	if err != nil {
		Error("geoip2 Error: " + err.Error())
	}
	defer db.Close()

	CountryCode := "NONE"

	ClientIP, ok := RequestJSON["client_ip"].(string)
	if !ok {
		ClientIP = "127.0.0.1"
	}

	ClientPort, ok := RequestJSON["client_port"].(string)
	if !ok {
		ClientPort = "8080"
	}

	ClientProxyData, ok := RequestJSON["proxy_data"].(string)
	if !ok {
		ClientProxyData = ""
	}

	ClientUserAgent, ok := RequestJSON["user_agent"].(string)
	if !ok {
		ClientUserAgent = "ISI"
	}

	Log("IP: " + ClientIP + " :: Country Code: " + CountryCode)

	params = make([]interface{}, 7)
	params[0] = InterfaceToInt(UserData["id"])
	params[1] = ClientIP
	params[2] = ClientPort
	params[3] = ClientProxyData
	params[4] = ClientUserAgent
	params[5] = CountryCode
	params[6] = AuthToken
	_, ok = UpdateDB("INSERT INTO `login_history`(`user_id`,`client_ip`,`client_port`,`proxy_data`,`user_agent`,`country_code`,`login_time`,`session_token`)VALUES(?,?,?,?,?,?,UNIX_TIMESTAMP(),?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 500
		ResponseJSON["message"] = GetMessageByID("en", 11) //"Server Error: Can't create session"
		return
	}

	ResponseJSON["auth_token"] = AuthToken
	ResponseJSON["user_id"] = InterfaceToInt(UserData["id"])
	ResponseJSON["username"] = UserData["username"]
	ResponseJSON["first_name"] = UserData["first_name"]
	ResponseJSON["last_name"] = UserData["last_name"]
	ResponseJSON["email"] = UserData["email"]
	ResponseJSON["phone_num"] = UserData["phone_num"]
	ResponseJSON["profile_pic"] = UserData["profile_pic"]
	ResponseJSON["created_datetime"] = UserData["created_datetime"]
	ResponseJSON["type"] = UserData["type"]

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 10) //"Login Successful"
}

// GetProfileInformation Get Profile Information
func GetProfileInformation(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `users` WHERE `id`=?"

	params := make([]interface{}, 1)
	params[0] = AuthData["user_id"]
	UserData, ok := GetSingleRow(Query, params, "default")
	if !ok {
		ResponseJSON["status"] = 401
		ResponseJSON["message"] = GetMessageByID("en", 31) //"Invalid Credentials"
		return
	}
	ResponseJSON["data"] = UserData
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListOpSockets List Op Sockets
func ListOpSockets(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT id, CONCAT(`proto`, ':', `host`, ':', `port`) as socket FROM `op_sockets`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddDispatcherGroup Add Dispatcher Group
func AddDispatcherGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	groupName, ok := RequestJSON["group_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 4)
		return
	}

	params := make([]interface{}, 1)
	params[0] = groupName
	_, ok = GetSingleRow("SELECT 'group_name'  FROM dispatcher_groups WHERE group_name=?", params, "default")
	if ok {
		ResponseJSON["status"] = 409
		ResponseJSON["message"] = GetMessageByID("en", 5)
		return
	}

	params = make([]interface{}, 1)
	params[0] = groupName
	_, ok = UpdateDB("INSERT INTO `dispatcher_groups` (`group_name`) VALUES (?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateDispatcherGroup Update Dispatcher Group
func UpdateDispatcherGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	groupName, ok := RequestJSON["group_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 4)
		return
	}

	params := make([]interface{}, 2)
	params[0] = groupName
	params[1] = id

	_, ok = UpdateDB("UPDATE `dispatcher_groups` SET `group_name`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DeleteDispatcherGroup Delete Dispatcher Group
func DeleteDispatcherGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `dispatcher_groups` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListDispatcherGroup List Dispatcher Group
func ListDispatcherGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `dispatcher_groups`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddAdmin Add Admin
func AddAdmin(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	gender, ok := RequestJSON["gender"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 32)
		return
	}

	firstName, ok := RequestJSON["first_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 33)
		return
	}

	lastName, ok := RequestJSON["last_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 34)
		return
	}

	username, ok := RequestJSON["username"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 35)
		return
	}
	params := make([]interface{}, 1)
	params[0] = username
	_, ok = GetSingleRow("SELECT 'id'  FROM users WHERE username=?", params, "default")
	if ok {
		ResponseJSON["status"] = 409
		ResponseJSON["message"] = GetMessageByID("en", 39)
		return
	}

	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 36)
		return
	}

	email, ok := RequestJSON["email"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 37)
		return
	}
	params = make([]interface{}, 1)
	params[0] = email
	_, ok = GetSingleRow("SELECT 'id'  FROM users WHERE email=?", params, "default")
	if ok {
		ResponseJSON["status"] = 409
		ResponseJSON["message"] = GetMessageByID("en", 40)
		return
	}

	phoneNum, ok := RequestJSON["phone_num"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 38)
		return
	}

	params = make([]interface{}, 7)
	params[0] = gender
	params[1] = firstName
	params[2] = lastName
	params[3] = username
	params[4] = password
	params[5] = email
	params[6] = phoneNum
	_, ok = UpdateDB("INSERT INTO `users` (`gender`,`first_name`,`last_name`,`username`,`password`,`email`,`phone_num`,`created_datetime`) VALUES (?,?,?,?,?,?,?,UNIX_TIMESTAMP())", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)
}

// UpdateAdmin Update Admin
func UpdateAdmin(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	gender, ok := RequestJSON["gender"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 32)
		return
	}

	firstName, ok := RequestJSON["first_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 33)
		return
	}

	lastName, ok := RequestJSON["last_name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 34)
		return
	}

	username, ok := RequestJSON["username"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 35)
		return
	}
	params := make([]interface{}, 2)
	params[0] = username
	params[1] = id
	_, ok = GetSingleRow("SELECT 'id'  FROM users WHERE username=? AND id!=?", params, "default")
	if ok {
		ResponseJSON["status"] = 409
		ResponseJSON["message"] = GetMessageByID("en", 39)
		return
	}

	email, ok := RequestJSON["email"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 37)
		return
	}
	params = make([]interface{}, 2)
	params[0] = email
	params[1] = id
	_, ok = GetSingleRow("SELECT 'id'  FROM users WHERE email=? AND id!=?", params, "default")
	if ok {
		ResponseJSON["status"] = 409
		ResponseJSON["message"] = GetMessageByID("en", 40)
		return
	}

	phoneNum, ok := RequestJSON["phone_num"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 38)
		return
	}

	params = make([]interface{}, 7)
	params[0] = gender
	params[1] = firstName
	params[2] = lastName
	params[3] = username
	params[4] = email
	params[5] = phoneNum
	params[6] = id

	_, ok = UpdateDB("UPDATE `users` SET `gender`=?, `first_name`=?, `last_name`=?, `username`=?, `email`=?, `phone_num`=?  WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DeleteAdmin Delete Admin
func DeleteAdmin(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `users` WHERE `id`=? AND `type`=2", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListAdmin List Admin
func ListAdmin(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `users` WHERE `type`=2"
	AdvancedSQLCondition := "AND"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddDispatcher Add Dispatcher
func AddDispatcher(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	dispatcherGroupId, ok := RequestJSON["data_center_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 6)
		return
	}

	destination, ok := RequestJSON["destination"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 7)
		return
	}

	socket := ""
	/*socket, ok := RequestJSON["socket"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 9)
		return
	}*/

	state, ok := RequestJSON["state"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 11)
		return
	}

	probeMode, ok := RequestJSON["probe_mode"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 18)
		return
	}

	weight, ok := RequestJSON["weight"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 19)
		return
	}

	priority, ok := RequestJSON["priority"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 20)
		return
	}

	attrs, ok := RequestJSON["attrs"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 21)
		return
	}

	description, ok := RequestJSON["description"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 22)
		return
	}

	params := make([]interface{}, 9)
	params[0] = dispatcherGroupId
	params[1] = destination
	params[2] = socket
	params[3] = state
	params[4] = probeMode
	params[5] = weight
	params[6] = priority
	params[7] = attrs
	params[8] = description

	_, ok = UpdateDB("INSERT INTO `dispatcher` (`setid`,`destination`,`socket`,`state`,`probe_mode`,`weight`,`priority`,`attrs`,`description`) VALUES (?,?,?,?,?,?,?,?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateDispatcher Update Dispatcher
func UpdateDispatcher(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	dispatcherGroupId, ok := RequestJSON["data_center_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 6)
		return
	}

	destination, ok := RequestJSON["destination"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 7)
		return
	}

	socket := ""
	/*socket, ok := RequestJSON["socket"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 9)
		return
	}*/

	state, ok := RequestJSON["state"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 11)
		return
	}

	probeMode, ok := RequestJSON["probe_mode"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 18)
		return
	}

	weight, ok := RequestJSON["weight"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 19)
		return
	}

	priority, ok := RequestJSON["priority"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 20)
		return
	}

	attrs, ok := RequestJSON["attrs"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 21)
		return
	}

	description, ok := RequestJSON["description"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 22)
		return
	}

	params := make([]interface{}, 10)
	params[0] = dispatcherGroupId
	params[1] = destination
	params[2] = socket
	params[3] = state
	params[4] = probeMode
	params[5] = weight
	params[6] = priority
	params[7] = attrs
	params[8] = description
	params[9] = id

	_, ok = UpdateDB("UPDATE `dispatcher` SET `setid`=?, `destination`=?, `socket`=?, `state`=?, `probe_mode`=?, `weight`=?, `priority`=?, `attrs`=?, `description`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DeleteDispatcher Delete Dispatcher
func DeleteDispatcher(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `dispatcher` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DisableGateway Disable Gateway
func DisableGateway(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 2)
	params[0] = "1"
	params[1] = id

	_, ok = UpdateDB("UPDATE `dispatcher` SET `status`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// EnableGateway Disable Gateway
func EnableGateway(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 2)
	params[0] = "0"
	params[1] = id

	_, ok = UpdateDB("UPDATE `dispatcher` SET `status`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListDispatcher List Dispatcher
func ListDispatcher(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT *, `setid` as data_center_id, (SELECT `group_name` FROM `dispatcher_groups` WHERE `id`=`dispatcher`.`setid`) as data_center_title FROM `dispatcher`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddSubscriberGroup Add Subscriber Group
func AddSubscriberGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	subscriberId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	dispatcherGroupId, ok := RequestJSON["data_center_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 6)
		return
	}

	params := make([]interface{}, 2)
	params[0] = subscriberId
	params[1] = dispatcherGroupId

	_, ok = UpdateDB("INSERT INTO `location_routing` (`location_id`,`data_center_id`) VALUES (?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateSubscriberGroup Update Subscriber Group
func UpdateSubscriberGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	subscriberId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	dispatcherGroupId, ok := RequestJSON["data_center_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 6)
		return
	}

	params := make([]interface{}, 3)
	params[0] = subscriberId
	params[1] = dispatcherGroupId
	params[2] = id

	_, ok = UpdateDB("UPDATE `location_routing` SET `location_id`=?, `data_center_id`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// RemoveRoute Remove Route
func RemoveRoute(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `location_routing` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	params = make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `subscriber_groups` WHERE `subscriber_id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListSubscriberGroup List Subscriber Group
func ListSubscriberGroup(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `location_routing`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddSubscriber Add Subscriber
func AddSubscriber(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	username, ok := RequestJSON["username"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 12)
		return
	}

	domain, ok := RequestJSON["domain"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 13)
		return
	}

	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 14)
		return
	}

	email, ok := RequestJSON["email"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 15)
		return
	}

	dispatcherGroupId, ok := RequestJSON["dispatcher_group_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 17)
		return
	}

	PasswordBytes := md5.Sum([]byte(InterfaceToString(password)))
	md5Password := fmt.Sprintf("%x", PasswordBytes)

	sha256Bytes := sha256.Sum256([]byte(InterfaceToString(password)))
	sha256Password := fmt.Sprintf("%x", sha256Bytes)

	sha512Bytes := sha512.Sum512([]byte(InterfaceToString(password)))
	sha512Password := fmt.Sprintf("%x", sha512Bytes)

	params := make([]interface{}, 7)
	params[0] = username
	params[1] = domain
	params[2] = md5Password
	params[3] = email
	params[4] = md5Password
	params[5] = sha256Password
	params[6] = sha512Password

	subscriberId, ok := UpdateDB("INSERT INTO `subscriber` (`username`,`domain`,`password`,`email_address`,`ha1`,`ha1_sha256`,`ha1_sha512t256`) VALUES (?,?,?,?,?,?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	params = make([]interface{}, 2)
	params[0] = subscriberId
	params[1] = dispatcherGroupId

	_, ok = UpdateDB("INSERT INTO `subscriber_groups` (`subscriber_id`,`dispatcher_group_id`) VALUES (?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateSubscriber Update Subscriber
func UpdateSubscriber(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	username, ok := RequestJSON["username"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 12)
		return
	}

	domain, ok := RequestJSON["domain"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 13)
		return
	}

	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 14)
		return
	}

	email, ok := RequestJSON["email"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 15)
		return
	}

	dispatcherGroupId, ok := RequestJSON["dispatcher_group_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 17)
		return
	}

	params := make([]interface{}, 5)
	params[0] = username
	params[1] = domain
	params[2] = password
	params[3] = email
	params[4] = id

	_, ok = UpdateDB("UPDATE `subscriber` SET `username`=?, `domain`=?, `password`=?, `email_address`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	params = make([]interface{}, 2)
	params[0] = dispatcherGroupId
	params[1] = id

	_, ok = UpdateDB("UPDATE `subscriber_groups` SET `dispatcher_group_id`=? WHERE `subscriber_id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DeleteSubscriber Delete Subscriber
func DeleteSubscriber(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `subscriber` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	params = make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `subscriber_groups` WHERE `subscriber_id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListSubscriber List Subscriber
func ListSubscriber(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT *, (SELECT `subscriber_groups`.`dispatcher_group_id` FROM `subscriber_groups` WHERE `subscriber_groups`.`subscriber_id`= `subscriber`.`id` LIMIT 1) as dispatcher_group_id, (SELECT `dispatcher_groups`.`group_name` FROM `dispatcher_groups` WHERE `dispatcher_groups`.`id`=(SELECT `subscriber_groups`.`dispatcher_group_id` FROM `subscriber_groups` WHERE `subscriber_groups`.`subscriber_id`= `subscriber`.`id` LIMIT 1)) as dispatcher_group_title FROM `subscriber`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// ListPermission List Subscriber
func ListPermission(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT context_info, id, ip, mask, pattern, port, proto, grp as location_id FROM `permissions`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddPermission Add Permission
func AddPermission(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	dispatcherGroupId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	ip, ok := RequestJSON["ip"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 24)
		return
	}

	mask, ok := RequestJSON["mask"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 25)
		return
	}

	port, ok := RequestJSON["port"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 26)
		return
	}

	proto, ok := RequestJSON["proto"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 27)
		return
	}

	pattern, ok := RequestJSON["pattern"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 28)
		return
	}

	contextInfo, ok := RequestJSON["context_info"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 29)
		return
	}

	params := make([]interface{}, 7)
	params[0] = dispatcherGroupId
	params[1] = ip
	params[2] = mask
	params[3] = port
	params[4] = proto
	params[5] = pattern
	params[6] = contextInfo

	_, ok = UpdateDB("INSERT INTO `permissions` (`grp`,`ip`,`mask`,`port`,`proto`,`pattern`,`context_info`) VALUES (?,?,?,?,?,?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdatePermission Update Permission
func UpdatePermission(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	dispatcherGroupId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	ip, ok := RequestJSON["ip"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 24)
		return
	}

	mask, ok := RequestJSON["mask"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 25)
		return
	}

	port, ok := RequestJSON["port"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 26)
		return
	}

	proto, ok := RequestJSON["proto"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 27)
		return
	}

	pattern, ok := RequestJSON["pattern"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 28)
		return
	}

	contextInfo, ok := RequestJSON["context_info"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 29)
		return
	}

	params := make([]interface{}, 8)
	params[0] = dispatcherGroupId
	params[1] = ip
	params[2] = mask
	params[3] = port
	params[4] = proto
	params[5] = pattern
	params[6] = contextInfo
	params[7] = id

	_, ok = UpdateDB("UPDATE `permissions` SET `grp`=?, `ip`=?, `mask`=?, `port`=?, `proto`=?, `pattern`=?, `context_info`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DeletePermission Delete Permission
func DeletePermission(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `permissions` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// GetLocationAnis Get Location Anis
func GetLocationAnis(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `location_anis`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddLocationAni Add Location Ani
func AddLocationAni(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	aniPrefix, ok := RequestJSON["ani_prefix"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 45)
		return
	} else {
		params := make([]interface{}, 1)
		_, ok := GetSingleRow("SELECT `ani_prefix` FROM `location_anis` WHERE `ani_prefix`= ?", params, "default")
		if ok {
			ResponseJSON["status"] = 400
			ResponseJSON["message"] = GetMessageByID("en", 46)
			return
		}
	}

	locationId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	params := make([]interface{}, 2)
	params[0] = locationId
	params[1] = aniPrefix

	_, ok = UpdateDB("INSERT INTO `location_anis` (`location_id`,`ani_prefix`) VALUES (?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateLocationAni Update Location Ani
func UpdateLocationAni(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	locationId, ok := RequestJSON["location_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 43)
		return
	}

	aniPrefix, ok := RequestJSON["ani_prefix"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 45)
		return
	}

	params := make([]interface{}, 3)
	params[0] = locationId
	params[1] = aniPrefix
	params[2] = id

	_, ok = UpdateDB("UPDATE `location_anis` SET `location_id`=?, `ani_prefix`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// RemoveLocationAni Remove Location Ani
func RemoveLocationAni(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `location_anis` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddCarrier Add Carrier
func AddCarrier(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	name, ok := RequestJSON["name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	params := make([]interface{}, 1)
	params[0] = name

	_, ok = UpdateDB("INSERT INTO `carriers` (`name`) VALUES (?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateCarrier Update Carrier
func UpdateCarrier(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	name, ok := RequestJSON["name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	params := make([]interface{}, 2)
	params[0] = name
	params[1] = id

	_, ok = UpdateDB("UPDATE `carriers` SET `name`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// RemoveCarrier Remove Carrier
func RemoveCarrier(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `carriers` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// DisableCarrier Disable Carrier
func DisableCarrier(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 2)
	params[0] = "inActive"
	params[1] = id

	_, ok = UpdateDB("UPDATE `carriers` SET `status`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// EnableCarrier Disable Carrier
func EnableCarrier(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 2)
	params[0] = "active"
	params[1] = id

	_, ok = UpdateDB("UPDATE `carriers` SET `status`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// GetCarriers Get Carriers
func GetCarriers(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `carriers`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddCarrierIp Add Carrier Ip
func AddCarrierIp(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	carrierId, ok := RequestJSON["carrier_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 44)
		return
	}

	description, ok := RequestJSON["description"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	carrierGatewayIp, ok := RequestJSON["carrier_gateway_ip"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 42)
		return
	}

	port, ok := RequestJSON["port"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 26)
		return
	}

	params := make([]interface{}, 3)
	params[0] = carrierId
	params[1] = description
	params[2] = InterfaceToString(carrierGatewayIp) + ":" + InterfaceToString(port)

	_, ok = UpdateDB("INSERT INTO `carrier_ips` (`group_id`,`description`,`dst_uri`) VALUES (?,?,?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateCarrierIp Update Carrier Ip
func UpdateCarrierIp(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	carrierId, ok := RequestJSON["carrier_id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 44)
		return
	}

	description, ok := RequestJSON["description"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	carrierGatewayIp, ok := RequestJSON["carrier_gateway_ip"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 42)
		return
	}

	port, ok := RequestJSON["port"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 26)
		return
	}

	params := make([]interface{}, 4)
	params[0] = carrierId
	params[1] = description
	params[2] = InterfaceToString(carrierGatewayIp) + ":" + InterfaceToString(port)
	params[3] = id

	_, ok = UpdateDB("UPDATE `carrier_ips` SET `group_id`=?, `description`=?, `dst_uri`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// RemoveCarrierIp Remove Carrier Ip
func RemoveCarrierIp(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `carrier_ips` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// GetCarrierIps Get Carrier Ips
func GetCarrierIps(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT description, dst_uri, id, probe_mode, resources, group_id as carrier_id FROM `carrier_ips`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// AddLocation Add Location
func AddLocation(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	name, ok := RequestJSON["name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	params := make([]interface{}, 1)
	params[0] = name

	_, ok = UpdateDB("INSERT INTO `locations` (`name`) VALUES (?)", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// UpdateLocation Update Location
func UpdateLocation(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	name, ok := RequestJSON["name"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 41)
		return
	}

	params := make([]interface{}, 2)
	params[0] = name
	params[1] = id

	_, ok = UpdateDB("UPDATE `locations` SET `name`=? WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// RemoveLocation Remove Location
func RemoveLocation(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	id, ok := RequestJSON["id"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = GetMessageByID("en", 30)
		return
	}

	params := make([]interface{}, 1)
	params[0] = id
	_, ok = UpdateDB("DELETE FROM `locations` WHERE `id`=?", params, "default")
	if !ok {
		ResponseJSON["status"] = 486
		ResponseJSON["message"] = GetMessageByID("en", 1)
		return
	}

	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// GetLocations Get Locations
func GetLocations(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}, AuthData map[string]string) {

	Query := "SELECT * FROM `locations`"
	AdvancedSQLCondition := "WHERE"

	StrSQL, ok := AdvancedSQL(RequestJSON, RequestJSON, AdvancedSQLCondition)
	if !ok {
		return
	}

	params := make([]interface{}, 0)
	list, _ := GetAllRows(Query+StrSQL, params, "default")
	if list == nil {
		ResponseJSON["list"] = make([]string, 0)
	} else {
		ResponseJSON["list"] = list
	}
	ResponseJSON["status"] = 200
	ResponseJSON["message"] = GetMessageByID("en", 0)

}

// GetUserIdentification identify user identification field's name & value
func GetUserIdentification(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}) (bool, string, string) {

	userid, ok := RequestJSON["username"]
	if ok {
		return true, "`username`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] == 405 {
		return false, "", ""
	}
	userid, ok = RequestJSON["phone_num"]
	if ok {
		return true, "`phone_num`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] == 405 {
		return false, "", ""
	}
	userid, ok = RequestJSON["email"]
	if ok {
		return true, "`email`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] != 405 {
		ResponseJSON["status"] = 406
		ResponseJSON["message"] = GetMessageByID("en", 7) //"No user identification available"
	}
	return false, "", ""
}
