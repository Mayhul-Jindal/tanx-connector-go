package client

import (
	"context"
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload string `json:"payload"`
}

/*
To get started with integrating tanX for your Dapp, we suggest you test your connectivity to our REST APIs.
Health endpoint can be used to check the status of the API.
*/
func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.healthURL.String(), nil)
	if err != nil {
		return HealthResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return HealthResponse{}, err
	}
	defer resp.Body.Close()

	var healthResponse HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	if err != nil {
		return HealthResponse{}, err
	}

	return healthResponse, nil
}
