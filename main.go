package main

import (
    "fmt"
    "net/http"
	"github.com/joho/godotenv"
    "github.com/o-klepatskyi/exchange-rate-notifier/ratefetcher"
    "github.com/o-klepatskyi/exchange-rate-notifier/mailsender"
)

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	mailsender.SubscribeEmail(email)
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
    rate := ratefetcher.GetCachedRate()
    if rate == 0 {
        http.Error(w, "", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "%.2f", rate)
}

func sendEmailsHandler(w http.ResponseWriter, r *http.Request) {
    rate := ratefetcher.GetCachedRate()
    if rate == 0 {
        fmt.Printf("Rate is not available yet")
        return
    }
    mailsender.SendEmails(rate)
}

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Warn: .env file was not loaded. If you are running inside docker, ingore this.")
    }

    go ratefetcher.RateFetchLoop() // Start the background rate fetcher

    http.HandleFunc("/rate", rateHandler)
    http.HandleFunc("/sendEmails", sendEmailsHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
