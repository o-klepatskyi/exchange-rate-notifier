package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type CurrencyRate struct {
    CurrencyCodeA int     `json:"currencyCodeA"`
    CurrencyCodeB int     `json:"currencyCodeB"`
    RateBuy       float64 `json:"rateBuy,omitempty"`
    RateSell      float64 `json:"rateSell,omitempty"`
    RateCross     float64 `json:"rateCross,omitempty"`
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
    client := &http.Client{
        Timeout: 1 * time.Second,
    }
    resp, err := client.Get("https://api.monobank.ua/bank/currency")
    if err != nil {
        fmt.Println("Error fetching data", err.Error())
        http.Error(w, "", http.StatusBadRequest)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("Error fetching data", resp.StatusCode)
        http.Error(w, "", http.StatusBadRequest)
        return
    }

    var rates []CurrencyRate
    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        fmt.Println("Error parsing data")
        http.Error(w, "", http.StatusBadRequest)
        return
    }

    for _, rate := range rates {
        if rate.CurrencyCodeA == 840 && rate.CurrencyCodeB == 980 {
            json.NewEncoder(w).Encode(rate.RateBuy)
            return
        }
    }
    fmt.Println("USD Rate not found")
    http.Error(w, "", http.StatusBadRequest)
}

func main() {
    http.HandleFunc("/rate", rateHandler)
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
