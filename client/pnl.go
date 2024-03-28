package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "status": "success",
  "message": "Retrieval Successful",
  "payload": [
    {
      "currency": "btc",
      "pnl_currency": "usdc",
      "total_credit": "30.2311686",
      "total_debit": "8.6779552",
      "total_credit_value": "951376.399690608",
      "total_debit_value": "263030.247889309",
      "average_buy_price": "31470.05040653995757213302",
      "average_sell_price": "30310.16429876349211851197",
      "average_balance_price": "31446.4028484392969497",
      "total_balance_value": "677771.031254780024102"
    },
    {
      "currency": "eth",
      "pnl_currency": "usdc",
      "total_credit": "30.2311686",
      "total_debit": "8.6779552",
      "total_credit_value": "951376.399690608",
      "total_debit_value": "263030.247889309",
      "average_buy_price": "31470.05040653995757213302",
      "average_sell_price": "30310.16429876349211851197",
      "average_balance_price": "31446.4028484392969497",
      "total_balance_value": "677771.031254780024102"
    }
  ]
}
*/
type PNLPayload struct {
	Currency            string `json:"currency"`
	PnlCurrency         string `json:"pnl_currency"`
	TotalCredit         string `json:"total_credit"`
	TotalDebit          string `json:"total_debit"`
	TotalCreditValue    string `json:"total_credit_value"`
	TotalDebitValue     string `json:"total_debit_value"`
	AverageBuyPrice     string `json:"average_buy_price"`
	AverageSellPrice    string `json:"average_sell_price"`
	AverageBalancePrice string `json:"average_balance_price"`
	TotalBalanceValue   string `json:"total_balance_value"`
}

type PNLResponse struct {
	Status  Status       `json:"status"`
	Message string       `json:"message"`
	Payload []PNLPayload `json:"payload"`
}

func (c *Client) PNL(ctx context.Context) (PNLResponse, error) {
	err := c.CheckAuth()
	if err != nil {
		return PNLResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.pnlURL.String(), nil)
	if err != nil {
		return PNLResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.jwtToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PNLResponse{}, err
	}
	defer resp.Body.Close()

	var pnlResponse PNLResponse
	err = json.NewDecoder(resp.Body).Decode(&pnlResponse)
	if err != nil {
		return PNLResponse{}, err
	}

	
	if pnlResponse.Status == ERROR {
		// handling 4xx and 5xx errors
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return PNLResponse{}, &ErrClient{
				Status: resp.StatusCode,
				Err:    fmt.Errorf(pnlResponse.Message),
			}
		} else if resp.StatusCode >= 500 {
			return PNLResponse{}, &ErrServer{
				Status: resp.StatusCode,
				Err:    fmt.Errorf(pnlResponse.Message),
			}
		}

		return PNLResponse{}, fmt.Errorf("status: %d\nerror: %s", resp.StatusCode, pnlResponse.Message)
	}
	


	return pnlResponse, nil
}
