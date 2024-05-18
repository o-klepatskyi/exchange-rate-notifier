package ratefetcher

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

var cachedRate float64

func fetchRate() {
	fmt.Println("Fetching exchange rate")
    client := &http.Client{
        Timeout: 1 * time.Second,
    }
    resp, err := client.Get("https://api.monobank.ua/bank/currency")
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error fetching data:", err, "status:", resp.StatusCode)
        return
    }
    defer resp.Body.Close()

    var rates []struct {
        CurrencyCodeA int     `json:"currencyCodeA"`
        CurrencyCodeB int     `json:"currencyCodeB"`
        RateBuy       float64 `json:"rateBuy,omitempty"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        fmt.Println("Error parsing data:", err)
        return
    }

    for _, rate := range rates {
        if rate.CurrencyCodeA == 840 && rate.CurrencyCodeB == 980 {
            cachedRate = rate.RateBuy
            return
        }
    }

    fmt.Println("Rate not found")
}

func RateFetchLoop() {
    fetchRate()
    ticker := time.NewTicker(10 * time.Second)
    for cachedRate == 0 {
        <-ticker.C
        fetchRate()
    }
    // Monobank updates its cached rate every 5 minutes, no need to do it more often
    ticker = time.NewTicker(5 * time.Minute)
    for {
        <-ticker.C
        fetchRate()
    }
}

func GetCachedRate() float64 {
    return cachedRate
}
