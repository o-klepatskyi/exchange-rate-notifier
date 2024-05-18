package main

import (
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/o-klepatskyi/exchange-rate-notifier/database"
	"github.com/o-klepatskyi/exchange-rate-notifier/mailsender"
	"github.com/o-klepatskyi/exchange-rate-notifier/ratefetcher"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Warn: .env file was not loaded. If you are running inside docker, ingore this.")
    }

	database.InitDB()
	database.CreateTable()

    go ratefetcher.RateFetchLoop() // Start the background rate fetcher

    http.HandleFunc("/rate", rateHandler)
    http.HandleFunc("/sendEmails", sendEmailsHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
