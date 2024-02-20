// https://github.com/weigj/go-odbc
package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"os"
	_ "strings"
)

var awsAvailable = false
var awsSession *session.Session
var AWSFromName string

//InitializeAWS ()
//func InitializeAWS(Credentials map[string]string) {
func InitializeAWS() {
	Log("--------------Initialize AWS-----------")
	awsAvailable = false
	var err error
	//awsAccessKey := ""
	//awsSecretKey := ""
	//awsRegion := ""

	value, ok := Config["AWS_ACCESS_KEY_ID"]
	if ok {
		os.Setenv("AWS_ACCESS_KEY_ID", value)
		//awsAccessKey = value
	} else {
		Error("Missing configuration AWS_ACCESS_KEY_ID")
		return
	}
	value, ok = Config["AWS_SECRET_ACCESS_KEY"]
	if ok {
		os.Setenv("AWS_SECRET_ACCESS_KEY", value)
		//awsSecretKey = value
	} else {
		Error("Missing configuration AWS_SECRET_ACCESS_KEY")
		return
	}
	value, ok = Config["AWS_DEFAULT_REGION"]
	if ok {
		os.Setenv("AWS_DEFAULT_REGION", value)
		//awsRegion = value
	} else {
		Error("Missing configuration AWS_DEFAULT_REGION")
		return
	}

	//os.Setenv("AWS_ACCESS_KEY_ID", Credentials["AWS_ACCESS_KEY_ID"])
	//os.Setenv("AWS_SECRET_ACCESS_KEY", Credentials["AWS_SECRET_ACCESS_KEY"])
	//os.Setenv("AWS_DEFAULT_REGION", Credentials["AWS_DEFAULT_REGION"])
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	//AWSFromName = Credentials["FROM"]
	AWSFromName = Config["AWS_FROM"]
	awsSession, err = session.NewSession()
	if err != nil {
		Error("InitializeAWS:: can't create AWS Session, Error:" + err.Error())
		return
	}
	awsAvailable = true
	return
}


// UploadFileToS3 saves a file to aws bucket and returns the url to // the file and an error if there's any
func UploadFileToS3(filePath string, fileName string) bool {
	// Open the file for use
	file, err := os.Open(filePath + fileName)
	if err != nil {
		Error("Failed to open file:: " + err.Error())
		return false
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	Response, err := s3.New(awsSession).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(Config["AWS_BUCKET"]),
		Key:           aws.String("/" + fileName),
		ACL:           aws.String("public-read"),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		//ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		Error("Error:: " + err.Error())
		return false
	}
	Log("Response:: " + fmt.Sprintf("%v", Response))
	_ = os.Remove(filePath + fileName)
	return true
}