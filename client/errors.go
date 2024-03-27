package client

import "fmt"

var (
	ErrMarketNotProvided     = fmt.Errorf("provide a market")
	ErrInvalidAmount         = fmt.Errorf("amount should be greater than 0")
	ErrNotLoggedIn           = fmt.Errorf("not logged in, please login first as this is a protected endpoint")
	ErrCoinNotFound          = fmt.Errorf("coin not found")
	ErrInvalidNetwork        = fmt.Errorf("invalid network")
	ErrInsufficientBalance   = fmt.Errorf("insufficient balance")
	ErrInsufficientAllowance = fmt.Errorf("insufficient allowance")
	ErrAllowanceTooLow       = fmt.Errorf("allowance too low")
	ErrAlreadyLoggedIn       = fmt.Errorf("already logged in")
)

// json decoding error
type ErrJSONDecoding struct {
	Err error
}

func (e *ErrJSONDecoding) Error() string {
	return fmt.Sprintf("error: %s\n", e.Err)
}

// 4xx error
type ErrClient struct {
	Status int
	Err error
}

func (e *ErrClient) Error() string {
	return fmt.Sprintf("status: %d\nerror: %s\n", e.Status, e.Err)
}

// 5xx error
type ErrServer struct {
	Status int
	Err error
}

func (e *ErrServer) Error() string {
	return fmt.Sprintf("status: %d\nerror: %s\n", e.Status, e.Err)
}



/*
{"status":"error","message":"Buying ETH is disabled in Testnet temporarily","payload":""}
*/
