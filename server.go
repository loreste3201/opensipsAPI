package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//Config map[string]string global application configuration map
var Config map[string]string

//LoadConfiguration API
func LoadConfiguration() bool {
	Info("********Loading Configuration*****")
	params := make([]interface{}, 0)
	//params[0] = 1
	ConfigData, ok := GetAllRows("SELECT * FROM `configurations`", params, "default")
	if !ok {
		Error("Can't load configuration from database")
		return false
	}
	//for _, ThisParam := range ConfigData {
	for i := 0; i < len(ConfigData); i++ {
		ThisParam := ConfigData[i]
		Debug(fmt.Sprintf("This Parameter Name: %s\tValue: %s", ThisParam["param"], ThisParam["value"]))
		Config[ThisParam["param"]] = ThisParam["value"]
	}

	Debug(fmt.Sprintf("Configurations : %v", Config))
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func main() {
	InitLogger()
	defer CloseLogger()
	Debug("Connecting DB")
	if !ConnectDB() {
		os.Exit(-1)
	}
	Debug("DB Connected")
	defer DisconnectDB()
	Info("==================================")
	params := make([]interface{}, 0)
	Row, ok := GetSingleRow("SELECT NOW() AS 'Time'", params, "default")
	if ok {
		Info("DB Time " + Row["Time"])
	}
	Config = make(map[string]string)
	if !LoadConfiguration() {
		os.Exit(-1)
	}
	Info("==================================")
	Log("Starting HTTP Server")
	go StartHTTPServer()
	time.Sleep(2 * time.Second)
	if !WebServiceStatus {
		Panic("Failed to start web server")
	}
	//InitializeDynamicAPIs()
	InitializeCache()
	//InitializeScheduler()
	//InitializeApplePush()
	InitializeLanguages()
	Log("******************************************")
	Log("*              Ready to Serve            *")
	Log("******************************************")
	var UserCommand = ""
	reader := bufio.NewReader(os.Stdin)
	for UserCommand != "exit" {
		fmt.Printf("VTS::CPE # ")
		UserCommand, _ = reader.ReadString('\n')
		UserCommand = strings.ToLower(strings.TrimSuffix(UserCommand, "\n"))
		if UserCommand != "" {
			fmt.Printf("\nCommand : " + UserCommand + "\n")
			switch UserCommand {
			case "reload":
				//LoadConfigurationFile()
				//LoadConfiguration()
				break
			case "exit":
				os.Exit(0)
			default:
				Log("Unknown Command")
				break
			}
		}
	}
}
