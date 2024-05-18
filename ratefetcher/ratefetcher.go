package ratefetcher

import (
    "encoding/json"
    "fmt"
    "net/http"
	"time"
)

var cachedRate float64

func FetchRate() bool {
	fmt.Println("Fetching exchange rate")
    client := &http.Client{
        Timeout: 1 * time.Second,
    }
    resp, err := client.Get("https://api.monobank.ua/bank/currency")
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error fetching data:", err, "status:", resp.StatusCode)
        return false
    }
    defer resp.Body.Close()

    var rates []struct {
        CurrencyCodeA int     `json:"currencyCodeA"`
        CurrencyCodeB int     `json:"currencyCodeB"`
        RateBuy       float64 `json:"rateBuy,omitempty"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        fmt.Println("Error parsing data:", err)
        return false
    }

    for _, rate := range rates {
        if rate.CurrencyCodeA == 840 && rate.CurrencyCodeB == 980 {
            cachedRate = rate.RateBuy
            return true
        }
    }

    fmt.Println("Rate not found")
	return false
}

func GetCachedRate() float64 {
    return cachedRate
}
