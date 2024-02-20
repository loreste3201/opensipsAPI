package main

import (
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	_ "io"
	"io/ioutil"
	"log"
	_ "math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	_ "reflect"
	"strconv"
	"strings"
	"time"
)

var appDomain = ""
var baseURL = ""
var httpPort = "80"
var maxUploadSize int64
var sessionLife int64
var basePath = ""
var htmlPath = "html"
var templatePath = "templates"
var filePath = "files"
var fileServerPrefix = "/files"
var uploadPath = "files"
var uploadPrefix = "/upload"
var s3uploadPrefix = "/s3upload"
var appAPIPrefix = "app"
var webAPIPrefix = "web"
var adminPrefix = "iAdmin"
var adminUser = "Admin"
var adminPass = "Admin"
var sslCertificate = ""
var sslCertificateKey = ""
var sslPort = "443"

// WebServiceStatus true/false the realtime status of web server
var WebServiceStatus = false

func initializeConfiguration() bool {
	value, ok := Config["APP_DOMAIN"]
	if ok {
		appDomain = value
	} else {
		Error("Missing configuration APP_DOMAIN")
		return false
	}
	value, ok = Config["BASE_URL"]
	if ok {
		baseURL = value
	} else {
		baseURL = "http://" + GetIP()
	}
	value, ok = Config["HTTP_PORT"]
	if ok {
		httpPort = value
	}
	value, ok = Config["MAX_UPLOAD"]
	if ok {
		maxUploadSize, _ = strconv.ParseInt(value, 10, 64)
	} else {
		maxUploadSize = 25 * 1024 * 1024
	}
	value, ok = Config["SESSION_LIFE"]
	if ok {
		sessionLife, _ = strconv.ParseInt(value, 10, 64)
	} else {
		sessionLife = 1 * 60 * 60
	}
	basePath = GetMyPath()
	value, ok = Config["HTML_PATH"]
	if ok {
		if value[0:1] == "/" {
			htmlPath = value
		} else {
			htmlPath = baseURL + "/" + value
		}
	} else {
		htmlPath = basePath + "/html"
	}
	value, ok = Config["TEMPLATE_PATH"]
	if ok {
		if value[0:1] == "/" {
			templatePath = value
		} else {
			templatePath = basePath + "/" + value
		}
	} else {
		templatePath = basePath + "/templates"
	}
	value, ok = Config["FILE_PATH"]
	if ok {
		if value[0:1] == "/" {
			filePath = value
		} else {
			filePath = basePath + "/" + value
		}
	} else {
		filePath = basePath + "/files"
	}
	value, ok = Config["FILE_SERVER_PREFIX"]
	if ok {
		fileServerPrefix = "/" + value
	}
	value, ok = Config["UPLOAD_PATH"]
	if ok {
		if value[0:1] == "/" {
			uploadPath = value
		} else {
			uploadPath = baseURL + "/" + value
		}
	} else {
		uploadPath = filePath
	}
	value, ok = Config["UPLOAD_SERVER_PREFIX"]
	if ok {
		uploadPrefix = "/" + value
	}
	value, ok = Config["S3_UPLOAD_SERVER_PREFIX"]
	if ok {
		s3uploadPrefix = "/" + value
	}
	value, ok = Config["APP_API_PREFIX"]
	if ok {
		appAPIPrefix = value
	}
	value, ok = Config["WEB_API_PREFIX"]
	if ok {
		webAPIPrefix = value
	}
	value, ok = Config["ADMIN_PREFIX"]
	if ok {
		adminPrefix = value
	}
	value, ok = Config["ADMIN_USER"]
	if ok {
		adminUser = value
	}
	value, ok = Config["ADMIN_PASS"]
	if ok {
		adminPass = value
	}
	value, ok = Config["SSL_certificate"]
	if ok {
		sslCertificate = value
	}
	value, ok = Config["SSL_certificate_key"]
	if ok {
		sslCertificateKey = value
	}
	value, ok = Config["SSL_PORT"]
	if ok {
		sslPort = value
	}

	return true
}

func controlCommand(r *http.Request) map[string]interface{} {
	ResponseJSON := make(map[string]interface{})
	//Get Request Body Data
	BodyText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseJSON["message"] = "Invalid Request Body"
		ResponseJSON["status"] = 406
		Warning("controlCommand: can't ready request body text, Error: " + err.Error())
		return ResponseJSON
	}
	Debug("Request Body Text:\t" + string(BodyText))
	RequestJSON := make(map[string]interface{})
	err = json.Unmarshal(BodyText, &RequestJSON)
	if err != nil {
		ResponseJSON["message"] = "Invalid JSON"
		ResponseJSON["status"] = 406
		Warning("controlCommand: Failed Parsing JSON: " + err.Error())
		return ResponseJSON
	}
	username, ok := RequestJSON["username"]
	if !ok {
		ResponseJSON["message"] = "Missing Credentials"
		ResponseJSON["status"] = 407
		Warning("controlCommand:: Missing username")
		return ResponseJSON
	}
	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["message"] = "Missing Credentials"
		ResponseJSON["status"] = 407
		Warning("controlCommand:: Missing passeord")
		return ResponseJSON
	}
	if (InterfaceToString(username) != adminUser) || (InterfaceToString(password) != adminPass) {
		ResponseJSON["message"] = "Invalid Credentials"
		ResponseJSON["status"] = 407
		Warning("controlCommand:: Invalid credentials")
		return ResponseJSON
	}
	command, ok := RequestJSON["command"]
	if !ok {
		ResponseJSON["message"] = "Missing command"
		ResponseJSON["status"] = 406
		Warning("controlCommand:: Missing command")
		return ResponseJSON
	}
	switch InterfaceToString(command) {
	case "reload api":
		InitializDynamicAPIs()
		ResponseJSON["message"] = "dynamic api reloaded"
		ResponseJSON["status"] = 200
		break
	/*case "reload apple":
	InitializeApplePush()
	ResponseJSON["message"] = "apple client reloaded"
	ResponseJSON["status"] = 200
	break*/
	/*case "reload aws":
	InitializeAWS()
	ResponseJSON["message"] = "aws client reloaded"
	ResponseJSON["status"] = 200
	break*/
	/*case "reload fcm":
	InitializeFCM()
	ResponseJSON["message"] = "fcm reloaded"
	ResponseJSON["status"] = 200
	break*/
	case "reload db":
		DisconnectDB()
		ResponseJSON["message"] = "db reloaded"
		if !ConnectDB() {
			ResponseJSON["message"] = "failed to connect DB"
		}
		ResponseJSON["status"] = 200
		break
	case "reload config":
		ResponseJSON["message"] = "config reloaded"
		if !LoadConfiguration() {
			ResponseJSON["message"] = "failed reloading config"
		}
		ResponseJSON["status"] = 200
		break
	case "reload cache":
		StopCache()
		InitializeCache()
		ResponseJSON["message"] = "cache reloaded"
		if !CacheAvailable {
			ResponseJSON["message"] = "failed reloading cache"
		}
		ResponseJSON["status"] = 200
		break
	case "exit":
		os.Exit(0)
		break
	default:
		ResponseJSON["message"] = "Invalid command"
		ResponseJSON["status"] = 404
		Warning("controlCommand:: unknown command :" + InterfaceToString(command))
		break
	}
	return ResponseJSON
}

// Authenticate (RequestJSON map[string]interface{}, ResJSON map[string]interface{}) Returns (map[string]string{}, bool)
func Authenticate(RequestJSON map[string]interface{}, ResJSON map[string]interface{}) (map[string]string, bool) {
	ok, AuthToken := GetKeyFromJSON(RequestJSON, ResJSON, "auth_token")
	if !ok {
		return nil, false
	}
	var ThisUserData map[string]string
	var AuthData map[string]string
	cacheData, status := CacheGet(AuthToken)
	if status == true {
		Debug("Authenticate:: Session found in cache")
		if err := json.Unmarshal(cacheData, &ThisUserData); err != nil {
			Error("Authenticate:: Failed parse session cache, seems crupted, deleting session cache")
			CacheDelete(AuthToken)
			return Authenticate(RequestJSON, RequestJSON)
		}
		ThisUserData["CurrentTime"] = strconv.FormatInt(CurrentTimeStamp(), 10)
	} else {
		Debug("Authenticate:: No Session cache, Fetching from DB")
		params := make([]interface{}, 1)
		params[0] = AuthToken
		AuthData, ok = GetSingleRow("SELECT `user_id`,`login_history`.*,UNIX_TIMESTAMP() AS CurrentTime ,`login_type` FROM `login_history` WHERE `session_token` = ? ORDER BY `id` DESC LIMIT 1", params, "default")
		if !ok {
			ResJSON["status"] = 407
			ResJSON["message"] = "authentication required"
			return nil, false
		}
		if AuthData["logout_time"] != "0" {
			ResJSON["status"] = 407
			ResJSON["message"] = "authentication required"
			return nil, false
		}

		UserTable := "users"
		if InterfaceToString(AuthData["login_type"]) == "WEB" {
			UserTable = "admin_users"
		}

		params = make([]interface{}, 1)
		params[0] = AuthData["user_id"]
		ThisUserData, _ = GetSingleRow("SELECT * FROM "+UserTable+" WHERE `id` = ?", params, "default")

		ThisUserData["user_id"] = AuthData["user_id"]
		ThisUserData["login_type"] = AuthData["login_type"]
		ThisUserData["session_token"] = AuthData["session_token"]
		ThisUserData["client_ip"] = AuthData["client_ip"]
		ThisUserData["country_code"] = AuthData["country_code"]

		if iBytes, err := json.Marshal(ThisUserData); err == nil {
			CacheSet(AuthToken, iBytes)
		} else {
			Error("Authenticate:: Failed to convert AuthData to JSON Error: " + err.Error())
		}
	}
	SessionStartTime, _ := strconv.ParseInt(AuthData["login_time"], 10, 64)
	CurrentTime, _ := strconv.ParseInt(AuthData["CurrentTime"], 10, 64)
	if (SessionStartTime + sessionLife) < CurrentTime {
		CacheDelete(AuthToken)
		ResJSON["status"] = 407
		ResJSON["message"] = "Sesion Expired"
		return nil, false
	}

	return ThisUserData, true
}

func getHTMLTemplate(TemplateName string, TemplateVars map[string]string) (string, bool) {
	ThisTemplate := templatePath + "/" + TemplateName + ".html"
	if _, err := os.Stat(ThisTemplate); os.IsNotExist(err) {
		Error("HTML Template File desn't exist at " + ThisTemplate)
		return "", false
	}
	content, err := ioutil.ReadFile(ThisTemplate)
	if err != nil {
		Error("Error reading HTML Template '" + ThisTemplate + "' Contents :" + err.Error())
		return "", false
	}
	TemplateData := strings.Replace(string(content), "{{base_url}}", baseURL, -1)
	for key, val := range TemplateVars {
		Debug(fmt.Sprintf("Template Var : %s \tValue : %s\n", key, val))
		TemplateData = strings.Replace(TemplateData, "{{"+key+"}}", val, -1)
	}
	return TemplateData, true
}

func addRestHeader(w http.ResponseWriter) {
	w.Header().Set("Server", "ToGee Web Server")
	w.Header().Set("DevAdmin", "imfanee@gmail.com")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
}

func processUploadRequest(RetValue map[string]interface{}, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		RetValue["message"] = err.Error()
		return
	}

	RequestJSON := make(map[string]interface{})
	RequestJSON["auth_token"] = r.PostFormValue("auth_token")
	AuthData, ok := Authenticate(RequestJSON, RetValue)
	if !ok {
		return
	}

	UploadType := r.PostFormValue("type")
	if UploadType == "" {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid Upload Type, "
		return
	}

	if !IsSQLSafe(UploadType) {
		RetValue["status"] = 407
		RetValue["message"] = fmt.Sprintf("prohibited chracters in type '%s'", UploadType)
		return
	}
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" &&
		filetype != "application/pdf" {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + filetype
		return
	}

	NenoSeconds := time.Now().UnixNano()
	NewfileName := strconv.FormatInt(NenoSeconds, 10)

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}
	NewfileName = NewfileName + fileEndings[0]
	AbasoluteFileName := filePath + NewfileName
	NewFile, err := os.Create(AbasoluteFileName)
	if err != nil {
		RetValue["status"] = 500
		RetValue["message"] = "can't upload, " + err.Error()
		return
	}

	defer NewFile.Close()
	if _, err := NewFile.Write(fileBytes); err != nil {
		RetValue["status"] = 500
		RetValue["message"] = "can't upload, " + err.Error()
		return
	}
	params := make([]interface{}, 3)
	params[0] = UploadType
	params[1] = NewfileName
	params[2] = AuthData["id"]
	Query := "INSERT INTO `uploaded_files`(`upload_type`,`file_name`,`uploaded_by`,`uploaded_time`)VALUES(?,?,?,NOW())"
	if _, ok := UpdateDB(Query, params, "default"); !ok {
		RetValue["status"] = 500
		RetValue["message"] = "Server Error: Can't update data"
		return
	}

	switch UploadType {
	case "profile_pic":
		params := make([]interface{}, 2)
		params[0] = NewfileName
		params[1] = AuthData["id"]
		Query := "UPDATE `application_users` SET `profile_pic`=? WHERE `id`=?"
		if _, ok := UpdateDB(Query, params, "default"); !ok {
			RetValue["status"] = 500
			RetValue["message"] = "Server Error: Can't update data"
			return
		}
		break
	default:
		RetValue["status"] = 200
		RetValue["file_name"] = NewfileName
		RetValue["message"] = "file uploaded"
		return
	}

	RetValue["status"] = 200
	RetValue["message"] = "Successful"
	RetValue["filename"] = NewfileName
}

func processS3UploadRequest(RetValue map[string]interface{}, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		RetValue["message"] = err.Error()
		return
	}

	RequestJSON := make(map[string]interface{})
	RequestJSON["auth_token"] = r.PostFormValue("auth_token")
	AuthData, ok := Authenticate(RequestJSON, RetValue)
	if !ok {
		return
	}

	UploadType := r.PostFormValue("type")
	if UploadType == "" {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid Upload Type, "
		return
	}

	if !IsSQLSafe(UploadType) {
		RetValue["status"] = 407
		RetValue["message"] = fmt.Sprintf("prohibited chracters in type '%s'", UploadType)
		return
	}
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" &&
		filetype != "application/pdf" && filetype != "video/mp4" {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + filetype
		return
	}

	NenoSeconds := time.Now().UnixNano()
	NewfileName := strconv.FormatInt(NenoSeconds, 10)

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RetValue["status"] = 488
		RetValue["message"] = "Invalid File, " + err.Error()
		return
	}

	NewfileName = NewfileName + fileEndings[0]
	AbasoluteFileName := filePath + NewfileName

	NewFile, err := os.Create(AbasoluteFileName)
	if err != nil {
		RetValue["status"] = 500
		RetValue["message"] = "can't upload, " + err.Error()
		return
	}

	defer NewFile.Close()
	if _, err := NewFile.Write(fileBytes); err != nil {
		RetValue["status"] = 500
		RetValue["message"] = "can't upload, " + err.Error()
		return
	}

	ok = UploadFileToS3("files/", NewfileName)
	if !ok {
		RetValue["status"] = 500
		RetValue["message"] = "can't upload, " + err.Error()
		return
	}

	params := make([]interface{}, 3)
	params[0] = UploadType
	params[1] = NewfileName
	params[2] = AuthData["id"]
	Query := "INSERT INTO `uploaded_files`(`upload_type`,`file_name`,`uploaded_by`,`uploaded_time`)VALUES(?,?,?,NOW())"
	if _, ok := UpdateDB(Query, params, "default"); !ok {
		RetValue["status"] = 500
		RetValue["message"] = "Server Error: Can't update data"
		return
	}
	RetValue["status"] = 200
	RetValue["message"] = "Successful"
	//RetValue["filename"] = Config["BASE_URL"] + "/download/" + NewfileName
	RetValue["filename"] = "https://s3." + Config["AWS_DEFAULT_REGION"] + ".amazonaws.com/" + Config["AWS_BUCKET"] + "/" + NewfileName
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	var RequestPath = r.URL.Path
	addRestHeader(w)
	Log("File Upload Request " + string(RequestPath))
	RetValue := make(map[string]interface{})
	RetValue["status"] = 486
	RetValue["message"] = "UnKnown"
	RetValue["data"] = "None"

	processUploadRequest(RetValue, r)

	Log(fmt.Sprintf("Request-Status: %v\t%v)", RetValue["status"], RetValue["message"]))
	jsonString, _ := json.Marshal(RetValue)
	fmt.Fprintf(w, "%v", string(jsonString))
	return
}

func s3FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	var RequestPath = r.URL.Path

	Log("S3 File Download Request " + string(RequestPath))

	fileName := strings.Split(RequestPath, "/")

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName[2]))
	w.Header().Set("Content-Type", "application/octet-stream")

	url := "https://" + Config["AWS_BUCKET"] + ".s3." + Config["AWS_DEFAULT_REGION"] + ".amazonaws.com/" + fileName[2]

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		Log(err.Error())
		return
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		Log(err.Error())
		return
	}
	return
}

func s3FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	var RequestPath = r.URL.Path
	addRestHeader(w)
	Log("File Upload Request " + string(RequestPath))
	RetValue := make(map[string]interface{})
	RetValue["status"] = 486
	RetValue["message"] = "UnKnown"
	RetValue["data"] = "None"

	processS3UploadRequest(RetValue, r)

	Log(fmt.Sprintf("Request-Status: %v\t%v)", RetValue["status"], RetValue["message"]))
	jsonString, _ := json.Marshal(RetValue)
	fmt.Fprintf(w, "%v", string(jsonString))
	return
}

func processAPI(ResponseJSON map[string]interface{}, RequestURLArray []string, r *http.Request) {
	Debug("API Request")
	if len(RequestURLArray) < 2 {
		ResponseJSON["message"] = "Incomplete API URI"
		Warning("processAPI: Missing API Type OR Name ")
		//http.Error(w, err.Error(), 500)
		return
	}
	//Initialize Request Data Map
	RequestJSON := make(map[string]interface{})
	RemoteConnection := strings.Split(r.RemoteAddr, ":")
	RequestJSON["client_ip"] = RemoteConnection[0]
	RequestJSON["client_port"] = RemoteConnection[1]
	RequestJSON["proxy_data"] = r.Header.Get("X-Forwarded-For")
	reqToken := r.Header.Get("Authorization")
	if reqToken != "" {
		splitToken := strings.Split(reqToken, "Bearer ")
		RequestJSON["auth_token"] = splitToken[1]
	}
	RequestJSON["user_agent"] = r.UserAgent()
	RequestJSON["request_method"] = r.Method
	//Get Request Body Data
	BodyText, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ResponseJSON["message"] = "Invalid Request Body"
		Warning("processAPI: can't ready request body text, Error: " + err.Error())
		//http.Error(w, err.Error(), 500)
		return
	}
	Debug("Request Body Text:\t" + string(BodyText))
	err = json.Unmarshal(BodyText, &RequestJSON)

	// Empty body exception

	//if err != nil {
	//	ResponseJSON["message"] = "Invalid JSON"
	//	Warning("processAPI: Failed Parsing JSON: " + err.Error())
	//	//http.Error(w, err.Error(), 500)
	//	return
	//}

	var endpoint = RequestURLArray[1]
	var endpointState = 1
	var apiVersion = "v0"
	if RequestURLArray[1] == "v1" || RequestURLArray[1] == "v2" || RequestURLArray[1] == "v0" {
		endpoint = RequestURLArray[2] + RequestURLArray[1]
		endpointState = 2
		apiVersion = RequestURLArray[1]
	}

	RequestJSON["api_version"] = apiVersion
	RequestJSON["endpoint"] = endpoint
	var AuthData map[string]string
	var Authenticated bool

	switch RequestURLArray[endpointState] {
	case "Login":
		break
	default:
		AuthData, Authenticated = Authenticate(RequestJSON, ResponseJSON)
		if !Authenticated {
			return
		}

		/*
			if RequestURLArray[1] == "web" {
				params := make([]interface{}, 2)
				params[0] = AuthData["user_id"]
				params[1] = RequestURLArray[2]
				_, ok := UpdateDB("INSERT INTO `admin_activity_log` (`admin_id`,`action`,`created_datetime`)VALUES(?,?,UNIX_TIMESTAMP())", params, "default")
				if !ok {
					Error("Can't insert activity data to DB")
					return
				}
			}
		*/
		break
	}

	RequestJSON["API_TYPE"] = strings.ToLower(RequestURLArray[0])

	API(RequestJSON, ResponseJSON, AuthData)
	//ResponseJSON["min_android_verion"] = Config["Min_Android_version"]
	//ResponseJSON["min_ios_verion"] = Config["Min_Ios_version"]
	return
}

func processWebRequest(RequestURLArray []string, RequestPath string, w http.ResponseWriter) {
	Log("HTTPS Contents Request")
	Debug(fmt.Sprintf("RequestURLArray :%v\tLength:\t%d,RequestPath:%s\t", RequestURLArray, len(RequestURLArray), RequestPath))
	ResponseData := "Welcome to my website!"
	RequestedPage := "index.html"
	if len(RequestURLArray) > 0 && RequestURLArray[len(RequestURLArray)-1] != "" {
		RequestedPage = RequestURLArray[len(RequestURLArray)-1]
	} else {
		RequestPath = RequestedPage
	}
	FileExtension := ""
	LastDot := strings.LastIndex(RequestedPage, ".")
	NameLength := len(RequestedPage)
	if (LastDot > 0) && (LastDot < NameLength) {
		FileExtension = string(RequestedPage[LastDot+1 : NameLength])
	} else {
		FileExtension = filepath.Ext(RequestedPage)
	}
	FileData, ok := ReadFileContents(htmlPath + "/" + RequestPath)
	if ok {
		ResponseData = FileData
	} else {
		TemplateVars := make(map[string]string)
		TempData, ok := getHTMLTemplate("404", TemplateVars)
		if ok {
			ResponseData = TempData
		}
	}

	Log("Requested File Extension : " + FileExtension)
	FileExtension = strings.Replace(FileExtension, ".", "", -1)
	Log("Requested File Extension : " + FileExtension)
	switch FileExtension {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "text/html")
	}
	//Log(ResponseData)
	w.Write([]byte(ResponseData))
}

func generalhttpRequestHandler(w http.ResponseWriter, r *http.Request) {
	Debug(r.Proto + " Request From: " + r.RemoteAddr + "\tURL: " + r.URL.Path)
	RequestPath := strings.Replace(r.URL.Path, "//", "/", -1)
	if len(RequestPath) > 1 {
		FirstChar := 0
		if string(RequestPath[0]) == "/" {
			FirstChar = 1
		}
		LastChar := len(RequestPath)
		if string(RequestPath[LastChar-1]) == "/" {
			LastChar = LastChar - 1
		}
		RequestPath = string(RequestPath[FirstChar:LastChar])
	}
	Debug(r.RemoteAddr + " => " + RequestPath)
	RequestURLArray := strings.Split(RequestPath, "/")
	if strings.ToLower(RequestURLArray[0]) == "api" || strings.ToLower(RequestURLArray[0]) == "adminapi" {
		addRestHeader(w)
		//Initialize Response Data Map
		RetValue := make(map[string]interface{})
		RetValue["status"] = 406
		RetValue["message"] = "Invalid API Request"
		requestIn := CurrentMilliSeconds()
		processAPI(RetValue, RequestURLArray, r)
		RetValue["process_time"] = CurrentMilliSeconds() - requestIn
		//we got back from API
		Debug(fmt.Sprintf("Request-Status: %d\t%s", RetValue["status"], RetValue["message"]))
		//Convert MAP to JSON
		jsonString, _ := json.Marshal(RetValue)
		Debug(fmt.Sprintf("\nResponse:: %v\n", string(jsonString)))
		//Send Response Back
		fmt.Fprintf(w, "%v", string(jsonString))
	} else if RequestURLArray[0] == adminPrefix {
		RetValue := controlCommand(r)
		//we got back from API
		Debug(fmt.Sprintf("Request-Status: %d\t%s", RetValue["status"], RetValue["message"]))
		//Convert MAP to JSON
		jsonString, _ := json.Marshal(RetValue)
		Debug(fmt.Sprintf("\nResponse:: %v\n", string(jsonString)))
		//Send Response Back
		fmt.Fprintf(w, "%v", string(jsonString))
	} else {
		Info("General Web Request")
		processWebRequest(RequestURLArray, RequestPath, w)
	}
	return
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

// StartHTTPServer start HTTP Service
func StartHTTPServer() {
	if !initializeConfiguration() {
		Error("Can't initilize HTTP configuration")
		return
	}
	WebServiceStatus = true
	http.HandleFunc("/", generalhttpRequestHandler)
	http.HandleFunc(uploadPrefix+"/", fileUploadHandler)
	http.HandleFunc(s3uploadPrefix+"/", s3FileUploadHandler)
	http.HandleFunc("/download/", s3FileDownloadHandler)

	fs := http.FileServer(http.Dir(filePath))
	http.Handle(fileServerPrefix+"/", http.StripPrefix(fileServerPrefix, fs))

	if Config["HTTP_FLAG"] == "1" {
		go http.ListenAndServe(":"+httpPort, nil)
	}

	if Config["HTTPS_FLAG"] == "1" {
		err := http.ListenAndServeTLS(":"+sslPort, sslCertificate, sslCertificateKey, nil)
		if err != nil {
			Log("Error: " + err.Error())
		}
	}
}
