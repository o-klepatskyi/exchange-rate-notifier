package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"github.com/o-klepatskyi/exchange-rate-notifier/database"
	"github.com/o-klepatskyi/exchange-rate-notifier/mailsender"
	"github.com/o-klepatskyi/exchange-rate-notifier/ratefetcher"
)

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "", http.StatusMethodNotAllowed)
        return
    }

	email := r.FormValue("email")
	mailsender.SubscribeEmail(email)
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
        http.Error(w, "", http.StatusMethodNotAllowed)
        return
    }

    rate := ratefetcher.GetCachedRate()
    if rate == 0 {
        http.Error(w, "", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "%.2f", rate)
}

func sendEmailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "", http.StatusMethodNotAllowed)
        return
    }

    rate := ratefetcher.GetCachedRate()
    if rate == 0 {
        fmt.Printf("Rate is not available yet")
        return
    }
    mailsender.SendEmails(rate)
}

func MailSendLoop() {
	// Send updates on subscribed emails every 24 hours
	for {
		rate := ratefetcher.GetCachedRate()
		if rate != 0 {
			mailsender.SendEmails(rate)
			time.Sleep(24 * time.Hour)
		} else {
			fmt.Println("Rate is invalid, retrying sending emails shortly...")
			time.Sleep(10 * time.Second)
		}
	}
}

func RateFetchLoop() {
    success := ratefetcher.FetchRate()

	// Make sure we fetch rate as fast as possible first time
    for !success {
        time.Sleep(10 * time.Second)
        success = ratefetcher.FetchRate()
    }

    // Monobank updates its cached rate every 5 minutes, no need to do it more often
    for {
        time.Sleep(5 * time.Minute)
        ratefetcher.FetchRate()
    }
}

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Warn: .env file was not loaded. If you are running inside docker, ingore this.")
    }

	database.InitDB()
	database.CreateTable()

    go RateFetchLoop()
	go MailSendLoop()

    http.HandleFunc("/rate", rateHandler)
    http.HandleFunc("/sendEmails", sendEmailsHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
