package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
)

var sysLog *syslog.Writer
var fileHandle *os.File

//var logToConsole = "noLog"
//var logToFile = "noLog"
//var logToSyslog = "noLog"

var theLogger = map[string]interface{}{
	"Error":   broadCastLog,
	"Warning": noLog,
	"Info":    noLog,
	"Debug":   noLog,
}

var broadCaster = map[string]interface{}{
	"logToConsole": noLog,
	"logToFile":    noLog,
	"logToSyslog":  noLog,
}

// CloseLogger close all log outputs
func CloseLogger() {
	broadCastLog("Shutting Down Logger")
	fileHandle.Close()
	sysLog.Close()
}

//InitLogger Initialize Logger to be used in allover the service
func InitLogger() {
	TempValue := os.Getenv("I_LOGCONSOLE")
	if (TempValue != "") && (strings.ToLower(TempValue) == "yes") {
		broadCaster["logToConsole"] = logConsole
	}

	TempValue = os.Getenv("I_LOGLFILE")
	if TempValue != "" {
		fileHandle, err := os.OpenFile(TempValue, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file:%s\t %v", TempValue, err)
		}
		//defer fileHandle.Close()
		log.SetOutput(fileHandle)
		//log.Println("Logger Initiated ")
		broadCaster["logToFile"] = toFile
	}
	var err error
	TempValue = os.Getenv("I_SYSLOG_URI")
	if TempValue != "" {
		SysLogParams := strings.Split(TempValue, ":")
		if (len(SysLogParams) < 3) || ((strings.ToLower(SysLogParams[0]) != "tcp") && (strings.ToLower(SysLogParams[0]) != "udp")) || (InterfaceToInt(SysLogParams[2]) < 1) {
			log.Fatalf("Invaid syslog URI :%s", TempValue)
		}
		sysLog, err = syslog.Dial(SysLogParams[0], SysLogParams[1]+":"+SysLogParams[2], syslog.LOG_WARNING|syslog.LOG_DAEMON, "iWeb")
	} else {
		sysLog, err = syslog.New(syslog.LOG_NOTICE, "iWeb")
	}
	if err != nil {
		log.Fatalf("Failed to connect syslog server :%s,\tError: %v", TempValue, err)
	}
	broadCaster["logToSyslog"] = toSyslog

	TempValue = os.Getenv("I_LOGLEVEL")
	switch TempValue {
	case "1", "Error":
		theLogger["Debug"] = noLog
		theLogger["Info"] = noLog
		theLogger["Warning"] = noLog
		theLogger["Error"] = broadCastLog
		break
	case "2", "Warning":
		theLogger["Debug"] = noLog
		theLogger["Info"] = noLog
		theLogger["Warning"] = broadCastLog
		theLogger["Error"] = broadCastLog
		break
	case "3", "Info":
		theLogger["Debug"] = noLog
		theLogger["Info"] = broadCastLog
		theLogger["Warning"] = broadCastLog
		theLogger["Error"] = broadCastLog
		break
	default:
		theLogger["Debug"] = broadCastLog
		theLogger["Info"] = broadCastLog
		theLogger["Warning"] = broadCastLog
		theLogger["Error"] = broadCastLog
	}
	broadCastLog("Logger Initialized")
}

/*func Log(Data string) {
	sysLog, err := syslog.Dial("tcp", "localhost:1234",
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "demotag")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(Data)
	sysLog.Notice(Data)
}*/

// Log to send logs to systemwide configured log outputs. call as Log(Data String)
func Log(Data string) {
	TxtColor := color.FgGreen.Render
	theLogger["Info"].(func(string))(TxtColor(Data))
}

// Error send logs to systemwide configured log outputs. call as Error(Data String)
func Error(Data string) {
	TxtColor := color.FgRed.Render
	theLogger["Error"].(func(string))(TxtColor(Data))
}

// Warning send logs to systemwide configured log outputs. call as Warning(Data String)
func Warning(Data string) {
	TxtColor := color.FgYellow.Render
	theLogger["Warning"].(func(string))(TxtColor(Data))
}

// Info send logs to systemwide configured log outputs. call as Info(Data String)
func Info(Data string) {
	TxtColor := color.FgBlue.Render
	theLogger["Info"].(func(string))(TxtColor(Data))
}

// Debug send logs to systemwide configured log outputs. call as Debug(Data String)
func Debug(Data string) {
	TxtColor := color.FgWhite.Render
	theLogger["Debug"].(func(string))(TxtColor(Data))
}

// Panic send logs to systemwide configured log outputs and exit. call as Panic(Data String)
func Panic(Data string) {
	TxtColor := color.FgLightRed.Render
	logConsole(TxtColor(Data))
	theLogger["Debug"].(func(string))(TxtColor(Data))
	os.Exit(-1)
}

func broadCastLog(Data string) {
	broadCaster["logToConsole"].(func(string))(Data)
	broadCaster["logToFile"].(func(string))(Data)
	broadCaster["logToSyslog"].(func(string))(Data)
}
func logConsole(Data string) {
	fmt.Printf("\n%s:\t%s\n", time.Now().String(), Data)
}

func toFile(Data string) {
	log.Println(Data)
}

func toSyslog(Data string) {
	sysLog.Notice(Data)
}

func noLog(Data string) {
	return
}
