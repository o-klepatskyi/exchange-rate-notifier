package mailsender

import (
	"fmt"
	"net/smtp"
	"net/mail"
	"os"
	"slices"
)

var emailList []string = []string{}

func isEmailValid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func SubscribeEmail(email string) bool {
	if !isEmailValid(email) {
		fmt.Println("Invalid email:", email)
		return false
	}
	if slices.Contains(emailList, email) {
		fmt.Println("Email already subscribed:", email)
		return false
	}
	emailList = append(emailList, email)
	fmt.Println("Email subscribed:", email)
	return true
}

func SendEmails(rate float64) {
	if len(emailList) == 0 {
		fmt.Println("No emails subscribed")
		return
	}
	smtpHost := os.Getenv("EXCHANGE_RATE_NOTIFIER_SMTP_HOST")
    smtpPort := os.Getenv("EXCHANGE_RATE_NOTIFIER_SMTP_PORT")
    smtpUser := os.Getenv("EXCHANGE_RATE_NOTIFIER_SMTP_USER")
    smtpPass := os.Getenv("EXCHANGE_RATE_NOTIFIER_SMTP_PASS")

    auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
    from := smtpUser
    subject := "Current USD Exchange Rate"
    body := fmt.Sprintf("The current exchange rate is %.2f UAH per USD.", rate)


	msg := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emailList, msg)
	if err != nil {
		fmt.Println("Error sending emails:", err.Error())
	} else {
		fmt.Println("Emails sent")
	}
}
