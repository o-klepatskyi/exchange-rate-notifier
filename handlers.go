package main

import (
	"fmt"
	"net/http"
	"github.com/o-klepatskyi/exchange-rate-notifier/mailsender"
	"github.com/o-klepatskyi/exchange-rate-notifier/ratefetcher"
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
