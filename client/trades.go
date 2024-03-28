package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderBy string

const (
	ASC  OrderBy = "asc"
	DESC OrderBy = "desc"
)

type TradesOptions struct {
	Limit     int
	Timestamp int64
	OrderBy   OrderBy
}

type TradesPayload struct {
	ID        int     `json:"id"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
	Total     float64 `json:"total"`
	Market    string  `json:"market"`
	CreatedAt int64   `json:"created_at"`
	TakerType string  `json:"taker_type"`
}

type TradesResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Payload []TradesPayload `json:"payload"`
}

/*
Retrieve the list of recent market trades placed successfully for a specific market using this endpoint
*/
func (c *Client) Trades(ctx context.Context, market string, options TradesOptions) (TradesResponse, error) {
	if len(market) == 0 {
		return TradesResponse{}, ErrMarketNotProvided
	}

	tradesURL := *c.tradesURL
	params := tradesURL.Query()

	params.Set("market", market)

	if options.Limit != 0 {
		params.Set("limit", fmt.Sprintf("%d", options.Limit))
	}

	if options.Timestamp != 0 {
		params.Set("timestamp", fmt.Sprintf("%d", options.Timestamp))
	}

	if options.OrderBy != "" {
		params.Set("from", string(options.OrderBy))
	}

	tradesURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tradesURL.String(), nil)
	if err != nil {
		return TradesResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TradesResponse{}, err
	}
	defer resp.Body.Close()

	var tradesResponse TradesResponse
	err = json.NewDecoder(resp.Body).Decode(&tradesResponse)
	if err != nil {
		return TradesResponse{}, err
	}

	return tradesResponse, nil
}
