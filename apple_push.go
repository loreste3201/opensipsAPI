package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"golang.org/x/net/http2"
)

var appleMessageJSON map[string]interface{}
var certificates tls.Certificate
var appleTopic string
var applePushAvailable = false
var pushType apns2.EPushType
var httpClientApple = &http.Client{}

//InitializeApplePush ()
func InitializeApplePush() {
	Log("InitializeApplePush")
	applePushAvailable = false
	var pushCertificateFile = ""
	var certificatePassword = ""
	var err error
	value, ok := Config["APPLE_PUSH_CERTIFICATE"]
	if ok {
		pushCertificateFile = value
	} else {
		Error("Missing configuration APPLE_PUSH_CERTIFICATE")
		return
	}
	value, ok = Config["APPLE_PUSH_CERTIFICATE_PASSWORD"]
	if ok {
		certificatePassword = value
	} else {
		Error("Missing configuration APPLE_PUSH_CERTIFICATE_PASSWORD")
		return
	}
	value, ok = Config["APPLE_APP_TOPIC"]
	if ok {
		appleTopic = value
	} else {
		Error("Missing configuration APPLE_APP_TOPIC")
		return
	}
	certificates, err = certificate.FromP12File(pushCertificateFile, certificatePassword)
	if err != nil {
		Log("InitializeApplePush:: Certificate Error:" + err.Error())
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificates},
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http2.Transport{TLSClientConfig: tlsConfig}
	httpClientApple = &http.Client{Transport: transport}

	appleMessageJSON = make(map[string]interface{})
	applePushAvailable = true
	return
}

//SendApplePush (UserID string, DeviceToken string, Title string, Data map[string]interface{}, DevMode bool)
func SendApplePush(UserID string, DeviceToken string, Title string, Data map[string]interface{}, DevMode bool, OsVer bool) bool {

	appleMessageJSON = make(map[string]interface{})
	logText := fmt.Sprintf("SendApplePush::\tTo:%s\tTitle:%s", DeviceToken, Title)
	Debug(logText)
	Debug("DevMode: " + fmt.Sprintf("%v", DevMode))
	params := make([]interface{}, 9)
	params[0] = "VOIP"
	params[1] = UserID
	params[2] = DeviceToken
	params[4] = "FAIL"
	if !applePushAvailable {
		params[6] = "Service not available"
		Error("SendApplePush:: Service not available")
		return false
	}
	aps := make(map[string]interface{})
	aps["alert"] = Title
	aps["sound"] = "default"
	aps["badge"] = 0
	if len(Data) > 0 {
		aps["content-available"] = 1
		aps["custom"] = Data
	} else {
		aps["content-available"] = 0
	}
	appleMessageJSON["aps"] = aps

	Debug(fmt.Sprintf("SendApplePush::%s", appleMessageJSON))

	jsonString, err := json.Marshal(appleMessageJSON)
	if err != nil {
		Error(logText + "\tcan't Marshal JSON, Error:" + err.Error())
		params[6] = "can't Marshal JSON, Error:" + err.Error()
		addApplePushCDR(params)
		return false
	}
	if DevMode {
		sendDevPush(DeviceToken, jsonString)
	}
	params[3] = jsonString

	pushType = apns2.PushTypeBackground

	notification := &apns2.Notification{}
	notification.DeviceToken = DeviceToken
	notification.Topic = appleTopic
	notification.Payload = jsonString
	notification.Priority = 10
	if OsVer {
		Log("Sending Push Type Header")
		notification.PushType = pushType
	}

	var client *apns2.Client
	if DevMode {
		client = apns2.NewClient(certificates).Development()
	} else {
		client = apns2.NewClient(certificates).Production()
	}
	Debug("Response: " + fmt.Sprintf("%v", notification))
	Response, err := client.Push(notification)
	Debug("Response: " + fmt.Sprintf("%v", Response))
	if err != nil {
		Error(logText + "\tError Sending push:" + err.Error())
		params[6] = "Error Sending push:" + err.Error()
		addApplePushCDR(params)
		return false
	}
	params[5] = Response.StatusCode
	params[6] = Response.Reason
	params[7] = Response.ApnsID
	params[8] = fmt.Sprintf("%v", Response)
	if Response.Sent() {
		params[4] = "SUCCESS"
		addApplePushCDR(params)
		return true
	}
	addApplePushCDR(params)
	return false
}

func sendDevPush(DeviceToken string, Payload []byte) {

	Log("------------------------Dev Push --------------------------")
	Log("Token: " + DeviceToken)
	Log("Payload: " + string(Payload))

	//body := strings.NewReader(`{"aps":{"alert":"audioCall from 923123811118","badge":0,"content-available":1,"custom":{"additional_param":"","call_type":"audio","created_at":1568287378,"from_user":"923123811118","from_voip_user":"USR923123811118","group_id":"331","group_title":"Test, Adil","notification_type":"join_call","temp_group_id":"1568287378433001"},"sound":"default"}}`)

	body := strings.NewReader(fmt.Sprintf("%v", string(Payload)))
	req, err := http.NewRequest("POST", "https://api.development.push.apple.com/3/device/"+DeviceToken, body)
	if err != nil {
		Log("Error:" + err.Error())
	}
	req.Header.Set("Apns-Topic", Config["APPLE_APP_TOPIC"])
	req.Header.Set("Apns-Priority", "10")
	req.Header.Set("Apns-Expiration", "0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClientApple.Do(req)
	if err != nil {
		Log("Error:" + err.Error())
	}
	Log(fmt.Sprintf("Response: %v", resp))
	defer resp.Body.Close()

}

func addApplePushCDR(params []interface{}) {
	UpdateDB("INSERT INTO `push_sent`(`push_type`,`user_id`,`device_token`, `sent_json`,`status`,`status_code`,`status_description`,`recieved_headers`,`recieved_body`,`created_datetime`)VALUES(?,?,?,?,?,?,?,?,?,UNIX_TIMESTAMP())", params, "default")
}
