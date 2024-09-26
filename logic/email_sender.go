package logic

import (
	"fmt"
	"medods_jwt_service/model"
	"net/smtp"
)

func SendEmail(user model.User, message string) {
	from := "medods@example.com"
	password := "medodsexample"

	to := []string{user.Email}
	smtpServer := model.SmtpServer{}
	smtpServer.Host = "smtp.gmail.com"
	smtpServer.Port = "587"

	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, []byte(message))

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("email sent")
}

func EmailTemplateMessage(user model.User) string {
	return fmt.Sprintf(
		"Dear, %s! Your IP has been changed to %s. If it wasnt you, contact us ASAP", user.Email, user.IP)
}
