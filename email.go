package main

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Email cutom data type to send emails
type Email struct {
	ToEmail   string
	ToName    string
	FromEmail string
	FromName  string
	Subject   string
	Body      string
	Carrier   string
	IsHTML    bool
	QueueID   int64
	UpdateDB  bool
}

// SMTPConf is Email sending configurations
type SMTPConf struct {
	AuthType            string
	ServerAddress       string
	ServerPort          int
	ServerConnection    string
	VerifiedCertificate bool
	Username            string
	Password            string
	SenderEmail         string
	SenderName          string
}

var emailQueue map[int64]Email
var emailAvailable = false
var smtpAvailable = false
var smtpCong SMTPConf
var smtpService = "SMTP"
var smtpConnection gomail.Dialer
var mailSender gomail.SendCloser
var sendMailBin string

var sendGridClient *sendgrid.Client

//InitializeEmail initialize and connect email sending service
func InitializeEmail() {
	var err error
	value, ok := Config["SMTP_SENDER_NAME"]
	if ok {
		smtpCong.SenderName = value
	} else {
		smtpCong.SenderName = Config["APP_DOMAIN"]
	}
	value, ok = Config["SMTP_SENDER_EMAIL"]
	if ok {
		smtpCong.SenderEmail = value
	} else {
		smtpCong.SenderEmail = "noreply@" + Config["APP_DOMAIN"]
	}
	value, ok = Config["SMTP_SERVER"]
	if ok {
		smtpCong.ServerAddress = value
	} else {
		Debug("No configuration SMTP_SERVER\nSearching local sendmail")
		connectLocalSendMail()
		return
	}
	value, ok = Config["SMTP_CONNECTION"]
	if ok {
		smtpCong.ServerConnection = value
	} else {
		smtpCong.ServerConnection = "PLAIN"
	}
	value, ok = Config["SMTP_PORT"]
	if ok {
		port, err := strconv.Atoi(value)
		if err != nil {
			Error("Invalid SMTP_PORT : " + value + "\t Error:" + err.Error())
			return
		}
		smtpCong.ServerPort = port
	} else {
		switch smtpCong.ServerConnection {
		case "STARTTLS":
			smtpCong.ServerPort = 465
			break
		case "TLS":
			smtpCong.ServerPort = 587
			break
		case "SSL":
			smtpCong.ServerPort = 465
			break
		default:
			smtpCong.ServerPort = 25
		}

	}
	if smtpCong.ServerConnection != "PLAIN" {
		smtpCong.VerifiedCertificate = true
		value, ok = Config["SMTP_INSECURE_CERTIFICATES"]
		if ok {
			if strings.ToLower(value) == "yes" {
				smtpCong.VerifiedCertificate = false
				smtpConnection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}
	value, ok = Config["SMTP_AUTH_TYPE"]
	if ok {
		smtpCong.AuthType = value
	} else {
		smtpConnection = gomail.Dialer{Host: smtpCong.ServerAddress, Port: smtpCong.ServerPort}
		if mailSender, err = smtpConnection.Dial(); err != nil {
			Error("can't connect SMTP  : Error: " + err.Error())
			Debug("Invalid SMTP configuration")
			//connectLocalSendMail()
			return
		}
		intitiateEmailer()
		return
	}
	value, ok = Config["SMTP_USERNAME"]
	if ok {
		smtpCong.Username = value
	} else {
		Error("Missing configuration SMTP_USERNAME")
		return
	}
	value, ok = Config["SMTP_PASSWORD"]
	if ok {
		smtpCong.Password = value
	} else {
		Error("Missing configuration SMTP_PASSWORD")
		return
	}
	/*
		if smtpCong.ServerConnection == "PLAIN" {
			smtpConnection = gomail.Dialer{Host: smtpCong.ServerAddress, Port: smtpCong.ServerPort, Username: smtpCong.Username, Password: smtpCong.Password}
		} else {

		}*/
	smtpConnection = gomail.Dialer{Host: smtpCong.ServerAddress, Port: smtpCong.ServerPort, Username: smtpCong.Username, Password: smtpCong.Password}
	if _, err = smtpConnection.Dial(); err != nil {
		Error("can't connect SMTP  : Error: " + err.Error())
		Debug("Invalid SMTP configuration\nUsing direct email sending option")
		//connectLocalSendMail()
		return
	}

	//Initialize Send Grid
	sendGridClient = sendgrid.NewSendClient(Config["SENDGRID_API_KEY"])
	intitiateEmailer()
	return

}

func connectLocalSendMail() {
	Debug("Trying to connect local sendmail")
	smtpConnection = gomail.Dialer{Host: "127.0.0.1", Port: 25}
	if _, err := smtpConnection.Dial(); err != nil {
		Error("can't connect sendmail" + err.Error())
		return
	}
	smtpService = "SENDMAIL"
	intitiateEmailer()
}

func intitiateEmailer() {
	emailQueue = make(map[int64]Email)
	emailAvailable = true
	processEmails()
}

// SendEmail to send email messages having params (ToEmail string, ToName string, subject string, message string, Msgtype string) and return type bool
func SendEmail(thisMail Email) (string, bool) {
	if !emailAvailable {
		return "No Email Sending Service", false
	}
	thisRegularExpression := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !thisRegularExpression.MatchString(thisMail.ToEmail) {
		return "Invalid ToEmail", false
	}
	params := make([]interface{}, 8)
	params[0] = thisMail.ToEmail
	if thisMail.ToName == "" {
		thisMail.ToName = thisMail.ToEmail
	}
	params[1] = thisMail.ToName
	if thisMail.FromEmail == "" {
		thisMail.FromEmail = smtpCong.SenderEmail
	}
	params[2] = thisMail.FromEmail
	if thisMail.FromName == "" {
		thisMail.FromName = smtpCong.SenderName
	}
	params[3] = thisMail.FromName
	if thisMail.Subject == "" {
		thisMail.Subject = "Email from " + thisMail.FromName
	}
	params[4] = thisMail.Subject
	if thisMail.Body == "" {
		return "Empty Email Body", false
	}
	params[5] = thisMail.Body
	params[6] = thisMail.IsHTML
	thisMail.QueueID = CurrentMilliSeconds()
	params[7] = thisMail.QueueID
	emailQueue[thisMail.QueueID] = thisMail
	lastID, status := UpdateDB("INSERT INTO `emails`(`to_email`,`to_user`,`from_email`,`from_user`,`subject`,`body`,`html`,`created_datetime`,`queue_id`)VALUES(?,?,?,?,?,?,?,UNIX_TIMESTAMP(),?)", params, "default")
	if status {
		thisMail.UpdateDB = true
		Debug("SendEmail:: Added to DB, ID: " + InterfaceToString(lastID))
	} else {
		thisMail.UpdateDB = false
		Error("SendEmail:: can't add to DB")
	}
	return "Added to Queue", true
}

func processEmails() {
	if !emailAvailable {
		return
	}
	timer := time.NewTimer(time.Second * 10)
	go func() {
		<-timer.C
		Debug("Processing email queue")
		for QueueID, thisMail := range emailQueue {
			params := make([]interface{}, 4)
			var statusText string
			var status bool
			if thisMail.Carrier == "Send-Grid" {
				params[0] = thisMail.Carrier

				statusText, status = sendMailWithSendGrid(thisMail)
				if status == true {
					delete(emailQueue, QueueID)
					params[1] = "SUCCESS"
				} else {
					params[1] = "FAILURE"
				}
			} else {
				params[0] = smtpService

				statusText, status = sendMail(thisMail)
				if status == true {
					delete(emailQueue, QueueID)
					params[1] = "SUCCESS"
				} else {
					params[1] = "FAILURE"
				}
			}

			params[2] = statusText
			params[3] = QueueID

			if updatedRows, status := UpdateDB("UPDATE `emails` SET `email_service`= ? , `status`= ? , `status_description`= ? WHERE `queue_id`= ?", params, "default"); status == true {
				Debug(fmt.Sprintf("processEmails:: Updated email status in DB, UpdatedRows: %d", updatedRows))
			} else {
				Error("processEmails:: Failed to updated DB")
			}
		}
		processEmails()
	}()
}

func sendMail(thisMail Email) (string, bool) {
	thisMessage := gomail.NewMessage()
	thisMessage.SetAddressHeader("From", thisMail.FromEmail, thisMail.FromName)
	thisMessage.SetAddressHeader("To", thisMail.ToEmail, thisMail.ToName)
	thisMessage.SetHeader("Subject", thisMail.Subject)
	if thisMail.IsHTML {
		thisMessage.SetBody("text/html", thisMail.Body)
	} else {
		thisMessage.SetBody("text/plain", thisMail.Body)
	}

	if err := smtpConnection.DialAndSend(thisMessage); err != nil {
		Error("sendSMTP:: can't send email, Error:" + err.Error())
		return err.Error(), false
	}
	return "Email Sent", true
}

func sendMailWithSendGrid(thisMail Email) (string, bool) {

	from := mail.NewEmail(thisMail.FromName, thisMail.FromEmail)
	subject := thisMail.Subject
	to := mail.NewEmail(thisMail.ToName, thisMail.ToEmail)

	plainTextContent := " "
	htmlContent := " "

	if thisMail.IsHTML {
		htmlContent = thisMail.Body
	} else {
		plainTextContent = thisMail.Body
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	response, err := sendGridClient.Send(message)
	if err != nil {
		Error("sendGrid:: can't send email, Error:" + err.Error())
		return err.Error(), false
	}

	if response.StatusCode != 202 {

		fmt.Println(response.StatusCode)
		Error("sendGrid:: can't send email, Error:" + response.Body)
		return response.Body, false
	}

	return "Email Sent", true

}
