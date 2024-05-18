package mailsender

import (
	"fmt"
	"net/smtp"
	"net/mail"
	"os"
	"github.com/o-klepatskyi/exchange-rate-notifier/database"
)

func isEmailValid(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func SubscribeEmail(email string) bool {
	if !isEmailValid(email) {
		fmt.Println("Invalid email:", email)
		return false
	}
	err := database.AddEmail(email)
	if err != nil {
		fmt.Println("Error adding email:", err)
		return false
	}
	fmt.Println("Email subscribed:", email)
	return true
}

func SendEmails(rate float64) {
	emailList, err := database.GetAllEmails()
	if err != nil {
		fmt.Println("Error fetching emails:", err)
		return
	}
	if len(emailList) == 0 {
		fmt.Println("No emails subscribed")
		return
	}
	smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USER")
    smtpPass := os.Getenv("SMTP_PASS")

    auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
    from := smtpUser
    subject := "Current USD Exchange Rate"
    body := fmt.Sprintf("The current exchange rate is %.2f UAH per USD.", rate)


	msg := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emailList, msg)
	if err != nil {
		fmt.Println("Error sending emails:", err.Error())
	} else {
		fmt.Printf("Emails sent to %d subscribers\n", len(emailList))
	}
}
