package mailsender

import (
	"fmt"
	"net/smtp"
	"os"
)

var emailList []string = []string{"test@gmail.com"}

func AddEmail(email string) {
	emailList = append(emailList, email)
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
	fmt.Println(smtpHost, smtpPort, smtpUser, smtpPass)

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
