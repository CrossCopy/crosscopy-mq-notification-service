package notification_service

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"sync"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

type EmailNotifier struct {
	User     string
	Password string
	Addr     string
	Host     string
	Sender   string
}

var lock = &sync.Mutex{}

var emailNotifier *EmailNotifier

func GetEmailNotifierInstance() *EmailNotifier {
	if emailNotifier == nil {
		lock.Lock()
		defer lock.Unlock()
		if emailNotifier == nil {
			fmt.Println("Creating EmailNotifier instance now.")
			emailNotifier = &EmailNotifier{}
		} else {
		}
	} else {
	}

	return emailNotifier
}

func (notifier *EmailNotifier) Init(email string, password string, serverAddr string, serverHost string) *EmailNotifier {
	notifier.User = email
	notifier.Password = password
	notifier.Addr = serverAddr
	notifier.Host = serverHost
	notifier.Sender = email
	return notifier
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}

func (notifier *EmailNotifier) SendSignupVerificationEmail(toEmail string, code string) *EmailNotifier {
	fmt.Printf("SendSignupVerificationEmail to %s\n", toEmail)
	sender := notifier.Sender
	to := []string{toEmail}
	subject := fmt.Sprintf("(%s) CrossCopy Signup Email Verification", code)
	bodyTemplate := `<!DOCTYPE html>
	<html lang="en">
	  <head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>CrossCopy Sign Up Verification</title>
		<style>
		  body {
			text-align: center
		  }
		</style>
	  </head>
	  <body>
		<h1>CrossCopy %s</h1>
		<p>Your verification code is <strong>%s</strong></p>
		<p>Please enter the verification code within 10 minutes</p>
	  </body>
	</html>
	`
	body := fmt.Sprintf(bodyTemplate, subject, code)
	request := Mail{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    body,
	}
	msg := BuildMessage(request)
	auth := smtp.PlainAuth("", notifier.User, notifier.Password, notifier.Host)
	err := smtp.SendMail(notifier.Addr, auth, sender, to, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	return notifier
}
