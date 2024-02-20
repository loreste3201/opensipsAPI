package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Languages map[string]interface{} Global Languages Map
var Languages map[string]map[int]interface{}

//InitializeLanguages from files
func InitializeLanguages() {

	Info("********Loading Language Files*****")

	var files []string

	root := "languages"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		Error(InterfaceToString(err))
	}

	Languages = make(map[string]map[int]interface{})
	for _, file := range files {

		jsonFile, err := os.Open(file)
		if err != nil {
			Error(InterfaceToString(err))
			return
		}
		Log("Successfully Opened " + file)

		byteValue, _ := ioutil.ReadAll(jsonFile)

		Language := make(map[int]interface{})
		json.Unmarshal(byteValue, &Language)

		fileName := strings.Split(file, "/")
		if len(fileName) > 1 {
			fileExtension := strings.Split(fileName[1], ".")
			Languages[fileExtension[0]] = Language
		}

		defer jsonFile.Close()
	}
}

//GetMessageByID func
func GetMessageByID(language string, id int) string {
	return InterfaceToString(Languages[language][id])
}
