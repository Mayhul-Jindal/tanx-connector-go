package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type BalancePayload struct {
	Currency       Currency `json:"currency"`
	Balance        string   `json:"balance"`
	Locked         string   `json:"locked"`
	DepositAddress string   `json:"deposit_address"`
}

type BalanceResponse struct {
	Status  Status         `json:"status"`
	Message string         `json:"message"`
	Payload BalancePayload `json:"payload"`
}

// Retrieve details of a specific user’s balance
func (c *Client) Balance(ctx context.Context, currency Currency) (BalanceResponse, error) {
	err := c.CheckAuth()
	if err != nil {
		return BalanceResponse{}, err
	}

	url := *c.balanceURL
	params := url.Query()
	params.Set("currency", string(currency))
	url.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return BalanceResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.jwtToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return BalanceResponse{}, err
	}
	defer resp.Body.Close()

	var balanceResponse BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&balanceResponse)
	if err != nil {
		return BalanceResponse{}, err
	}

	if balanceResponse.Status == ERROR {
		// handling 4xx and 5xx errors
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return BalanceResponse{}, &ErrClient{
				Status: resp.StatusCode,
				Err:    fmt.Errorf(balanceResponse.Message),
			}
		} else if resp.StatusCode >= 500 {
			return BalanceResponse{}, &ErrServer{
				Status: resp.StatusCode,
				Err:    fmt.Errorf(balanceResponse.Message),
			}
		}

		return BalanceResponse{}, fmt.Errorf("status: %d\nerror: %s", resp.StatusCode, balanceResponse.Message)
	}

	return balanceResponse, nil
}
