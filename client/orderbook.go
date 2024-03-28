package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderBookOptions struct {
	AskLimit int
	BidLimit int
}

type Ask struct {
	ID              int    `json:"id"`
	UUID            string `json:"uuid"`
	Side            string `json:"side"`
	OrdType         string `json:"ord_type"`
	Price           string `json:"price"`
	AvgPrice        string `json:"avg_price"`
	State           string `json:"state"`
	Market          string `json:"market"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	OriginVolume    string `json:"origin_volume"`
	RemainingVolume string `json:"remaining_volume"`
	ExecutedVolume  string `json:"executed_volume"`
	MakerFee        string `json:"maker_fee"`
	TakerFee        string `json:"taker_fee"`
	TradesCount     int    `json:"trades_count"`
}

type Bid struct {
	ID              int    `json:"id"`
	UUID            string `json:"uuid"`
	Side            string `json:"side"`
	OrdType         string `json:"ord_type"`
	Price           string `json:"price"`
	AvgPrice        string `json:"avg_price"`
	State           string `json:"state"`
	Market          string `json:"market"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	OriginVolume    string `json:"origin_volume"`
	RemainingVolume string `json:"remaining_volume"`
	ExecutedVolume  string `json:"executed_volume"`
	MakerFee        string `json:"maker_fee"`
	TakerFee        string `json:"taker_fee"`
	TradesCount     int    `json:"trades_count"`
}

type OrderBookPayload struct {
	Asks []Ask `json:"asks"`
	Bids []Bid `json:"bids"`
}

type OrderBookResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Payload OrderBookPayload `json:"payload"`
}

/*
Retrieve the latest asks and bids created within the tanX orderbook in a specific market using this endpoint
*/
func (c *Client) OrderBook(ctx context.Context, market string, options OrderBookOptions) (OrderBookResponse, error) {
	if len(market) == 0 {
		return OrderBookResponse{}, ErrMarketNotProvided
	}

	orderBookURL := *c.orderBookURL
	params := orderBookURL.Query()

	params.Set("market", market)

	if options.AskLimit != 0 {
		params.Set("ask_limit", fmt.Sprintf("%d", options.AskLimit))
	}

	if options.BidLimit != 0 {
		params.Set("bid_limit", fmt.Sprintf("%d", options.BidLimit))
	}

	orderBookURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, orderBookURL.String(), nil)
	if err != nil {
		return OrderBookResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OrderBookResponse{}, err
	}
	defer resp.Body.Close()

	var orderBookResponse OrderBookResponse
	err = json.NewDecoder(resp.Body).Decode(&orderBookResponse)
	if err != nil {
		return OrderBookResponse{}, err
	}

	return orderBookResponse, nil
}
