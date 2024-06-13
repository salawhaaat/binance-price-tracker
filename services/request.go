package services

import (
	"fmt"
	"net/url"
)

const (
	BaseURL      = "https://testnet.binance.vision"
	BaseEndpoint = "https://testnet.binance.vision/api/v3/ticker/price"
)

// prepareRequest prepares the request URL
func prepareRequest(symbol *string) string {
	reqURL, err := url.Parse(BaseEndpoint)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	params := url.Values{
		"symbol": []string{*symbol},
	}
	reqURL.RawQuery = params.Encode()
	return reqURL.String()
}
