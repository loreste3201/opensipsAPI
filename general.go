package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

//GetMyPath () retrun the absolute current folder path from where app is launched
func GetMyPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Debug(fmt.Sprintf("con't get MyPath, Error:%s", err.Error()))
		return ""
	}
	return dir
}

//GetIP () returns system's IP
func GetIP() string {
	Debug("GetIP: Get System's Local IP")
	ipAddress := ""
	ifaces, err := net.Interfaces()
	if err != nil {
		Error("GetIP: can't get interfaces, Error:" + err.Error())
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			Error("GetIP: can't get interface adress, Error:" + err.Error())
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil {
				strIP := fmt.Sprintf("%v", ip)
				Info("This Interface IP " + strIP)

				if strIP != "127.0.0.1" {
					ipAddress = strIP
				}
			}
		}
	}
	return ipAddress
}

//CurrentMilliSeconds () returns current millisecond time as int64
func CurrentMilliSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// CurrentTimeStamp () return current Unix Timestamp as int64
func CurrentTimeStamp() int64 {
	return int64(time.Now().Unix())
}

//isNumeric (string ) checl if provided string is numeric
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

//GetKeyFromJSON (RequestJSON map[string]interface{}, ResJSON map[string]interface{}, KeyName string) try to get value of key from json and return (bool status, string value)
func GetKeyFromJSON(RequestJSON map[string]interface{}, ResJSON map[string]interface{}, KeyName string) (bool, string) {
	value, ok := RequestJSON[KeyName]
	if !ok {
		ResJSON["status"] = 407
		ResJSON["message"] = KeyName + " required"
		return false, ""
	}
	KeyValue := InterfaceToString(value)
	if IsSQLSafe(KeyValue) {
		return true, KeyValue
	}
	ResJSON["status"] = 405
	ResJSON["message"] = "prohibited chracters in key " + KeyName
	return false, ""
}

//InterfaceToString (interface{}) get value of an interface return as string {
func InterfaceToString(ThisInterface interface{}) string {
	if ThisInterface == nil {
		Log("InterfaceToString :: NIL value passed for conversion")
		return ""
	}
	switch ThisInterface.(type) {
	case int:
		return strconv.Itoa(ThisInterface.(int))
	case int64:
		return strconv.FormatInt(ThisInterface.(int64), 10)
	case float64:
		TmpStr := strconv.FormatFloat(ThisInterface.(float64), 'f', 10, 64)
		if TmpStr[(len(TmpStr)-10):] == "0000000000" {
			return TmpStr[:(len(TmpStr) - 11)]
		}
		return TmpStr
	default:
		return ThisInterface.(string)
	}
}

//InterfaceToInt (interface{}) get value of interface and return as int64
func InterfaceToInt(ThisInterface interface{}) int64 {
	switch ThisInterface.(type) {
	case int:
		return ThisInterface.(int64)
	case int64:
		return ThisInterface.(int64)
	case float64:
		//return ThisInterface.(int64)
		return int64(ThisInterface.(float64))
	default:
		val, _ := strconv.ParseInt(ThisInterface.(string), 10, 64)
		return val
	}
}

//ReadFileContents (AbsoluteFilePath string) returns (contents as string, status as bool)
func ReadFileContents(FileName string) (string, bool) {
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		Log("File desn't exist at " + FileName)
		return err.Error(), false
	}
	content, err := ioutil.ReadFile(FileName)
	if err != nil {
		Log("Error reading File '" + FileName + "' Contents :" + err.Error())
		return err.Error(), false
	}
	//Log(string(content))
	return string(content), true
}

var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
