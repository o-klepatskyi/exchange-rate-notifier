package main

import (
    "fmt"
    "net/http"
    "github.com/o-klepatskyi/exchange-rate-notifier/ratefetcher"
)

func rateHandler(w http.ResponseWriter, r *http.Request) {
    rate := ratefetcher.GetCachedRate()
    if rate == 0 {
        http.Error(w, "", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "%.2f", rate)
}

func main() {
    go ratefetcher.StartRateFetcher() // Start the background rate fetcher

    http.HandleFunc("/rate", rateHandler)
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
