package client

import (
	"context"
	"encoding/json"
	"net/http"
)

type Ticker struct {
	Low                string `json:"low"`
	High               string `json:"high"`
	Open               string `json:"open"`
	Last               string `json:"last"`
	Volume             string `json:"volume"`
	Amount             string `json:"amount"`
	Vol                string `json:"vol"`
	AvgPrice           string `json:"avg_price"`
	PriceChangePercent string `json:"price_change_percent"`
	At                 string `json:"at"`
}

type TickerPayload struct {
	At     string `json:"at"`
	Ticker Ticker `json:"ticker"`
}

type TickerResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Payload TickerPayload `json:"payload"`
}

type AllTickerResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Payload map[string]TickerPayload `json:"payload"`
}

/*
Track 24 hour rolling window price changes to receive stats about a certain market.
For example 24 hour low, high, average price, price change and more.
*/
func (c *Client) Ticker(ctx context.Context, market string) (TickerResponse, error) {
	if len(market) == 0 {
		return TickerResponse{}, ErrMarketNotProvided
	}

	tickerURL := *c.tickerURL
	params := tickerURL.Query()
	params.Set("market", market)
	tickerURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tickerURL.String(), nil)
	if err != nil {
		return TickerResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TickerResponse{}, err
	}
	defer resp.Body.Close()

	var tickerResponse TickerResponse
	err = json.NewDecoder(resp.Body).Decode(&tickerResponse)
	if err != nil {
		return TickerResponse{}, err
	}

	return tickerResponse, nil
}

// Track 24 hour rolling window price changes to receive stats about all markets.
func (c *Client) AllTickers(ctx context.Context) (AllTickerResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.tickerURL.String(), nil)
	if err != nil {
		return AllTickerResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AllTickerResponse{}, err
	}
	defer resp.Body.Close()

	var allTickerResponse AllTickerResponse
	err = json.NewDecoder(resp.Body).Decode(&allTickerResponse)
	if err != nil {
		return AllTickerResponse{}, err
	}

	return allTickerResponse, nil
}
