package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OrderNonceRequest struct {
	Market  string  `json:"market"`
	OrdType string  `json:"ord_type"`
	Price   float64 `json:"price"`
	Side    string  `json:"side"`
	Volume  float64 `json:"volume"`
}

type OrderNonceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload struct {
		Nonce   int    `json:"nonce"`
		MsgHash string `json:"msg_hash"`
	} `json:"payload"`
}

func (c *Client) OrderNonce(ctx context.Context, market string, ordType string, price float64, side string, volume float64) (OrderNonceResponse, error) {
	orderNonceRequest := OrderNonceRequest{
		Market:  market,
		OrdType: ordType,
		Price:   price,
		Side:    side,
		Volume:  volume,
	}

	requestBody, err := json.Marshal(orderNonceRequest)
	if err != nil {
		return OrderNonceResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.orderNonceURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return OrderNonceResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.jwtToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OrderNonceResponse{}, err
	}
	defer resp.Body.Close()

	var orderNonceResponse OrderNonceResponse
	err = json.NewDecoder(resp.Body).Decode(&orderNonceResponse)
	if err != nil {
		return OrderNonceResponse{}, err
	}

	return orderNonceResponse, nil
}

type OrderCreateSignature struct {
	R string `json:"r"`
	S string `json:"s"`
}

type OrderCreateRequest struct {
	MsgHash   string               `json:"msg_hash"`
	Signature OrderCreateSignature `json:"signature"`
	Nonce     int                  `json:"nonce"`
}

type OrderCreateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload struct {
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
	} `json:"payload"`
}

func (c *Client) OrderCreate(ctx context.Context, starkPrivateKey string, market string, ordType string, price float64, side string, volume float64) (OrderCreateResponse, error) {
	orderNonceResponse, err := c.OrderNonce(ctx, market, ordType, price, side, volume)
	if err != nil {
		return OrderCreateResponse{}, err
	}
	log.Println("orderNonceResponse: ", orderNonceResponse)

	// todo
	/*
		1. The sign function is thrwoing runtime error
		2. Noticing the stack trace error
	*/
	something, err := Sign(starkPrivateKey, orderNonceResponse.Payload.MsgHash)
	if err != nil {
		return OrderCreateResponse{}, err
	}	
	log.Println(something)

	// todo
	/*
		1. This part canonly betouched when the  upper part is fixed
		1. use r and s values that comes after signing
	*/
	orderCreateRequest := OrderCreateRequest{
		MsgHash: orderNonceResponse.Payload.MsgHash,
		Signature: OrderCreateSignature{
			R: "",
			S: "",
		},
		Nonce: orderNonceResponse.Payload.Nonce,
	}

	requestBody, err := json.Marshal(orderCreateRequest)
	if err != nil {
		return OrderCreateResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.orderCreateURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return OrderCreateResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.jwtToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OrderCreateResponse{}, err
	}
	defer resp.Body.Close()

	var orderCreateResponse OrderCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&orderCreateResponse)
	if err != nil {
		return OrderCreateResponse{}, err
	}

	return orderCreateResponse, nil
}
