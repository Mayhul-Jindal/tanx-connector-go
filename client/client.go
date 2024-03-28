package client

import (
	"net/http"
	"net/url"
)

type Status string

const (
	SUCCESS Status = "success"
	ERROR   Status = "error"
)

// base url testnet/mainnet
type BaseURL string

const (
	MAINNET BaseURL = "https://api.tanx.fi"
	TESTNET BaseURL = "https://api-testnet.tanx.fi"
)

// endpoints
const (
	HEALTH_ENDPOINT      = "/sapi/v1/health/"
	TICKER_ENDPOINT      = "/sapi/v1/market/tickers/"
	CANDLESTICK_ENDPOINT = "/sapi/v1/market/kline/"
	ORDERBOOK_ENDPOINT   = "/sapi/v1/market/orderbook/"
	TRADES_ENDPOINT      = "/sapi/v1/market/trades/"

	NONCE_ENDPOINT         = "/sapi/v2/auth/nonce/"
	LOGIN_ENDPOINT         = "/sapi/v2/auth/login/"
	REFRESH_TOKEN_ENDPOINT = "/sapi/v1/auth/token/refresh/"
	PROFILE_ENDPOINT       = "/sapi/v1/user/profile/"
	BALANCE_ENDPOINT       = "/sapi/v1/user/balance/"
	PNL_ENDPOINT           = "/sapi/v1/user/pnl/"

	ORDER_BASE_ENDPOINT                = "/sapi/v1/orders/"
	ORDER_NONCE_ENDPOINT               = "/sapi/v1/orders/nonce/"
	ORDER_CREATE_ENDPOINT              = "/sapi/v1/orders/create/"
	ORDER_CANCEL_ENDPOINT              = "/sapi/v1/orders/cancel/"
	TRADES_LIST_ENDPOINT               = "/sapi/v1/trades/"	

	COIN_ENDPOINT                      = "/main/stat/v2/coins/"
	VAULTID_ENDPOINT                   = "/main/user/create_vault/"
	NETWORK_CONFIG_ENDPOINT            = "/main/stat/v2/app-and-markets/"
	CRYPTO_DEPOSIT_START_ENDPOINT      = "/sapi/v1/payment/stark/start/"
	CROSS_CHAIN_DEPOSIT_START_ENDPOINT = "/sapi/v1/deposits/crosschain/create/"
	LIST_DEPOSITS_ENDPOINT             = "/sapi/v1/deposits/all/"

	INTERNAL_TRANSFER_INITIATE_ENDPOINT = "/sapi/v1/internal_transfers/v2/initiate/"
	INTERNAL_TRANSFER_PROCESS_ENDPOINT  = "/sapi/v1/internal_transfers/v2/process/"
	INTERNAL_TRANSFER_GET_ENDPOINT      = "/sapi/v1/internal_transfers/v2/"
	INTERNAL_TRANSFER_USER_ENDPOINT     = "/sapi/v1/internal_transfers/v2/check_user_exists/"
	INTERNAL_TRANSFER_LIST_ENDPOINT     = "/sapi/v1/internal_transfers/v2/"

	START_NORMAL_WITHDRAWAL_ENDPOINT    = "/sapi/v1/payment/withdrawals/v1/initiate/"
	VALIDATE_NORMAL_WITHDRAWAL_ENDPOINT = "/sapi/v1/payment/withdrawals/v1/validate/"
	START_FAST_WITHDRAWAL_ENDPOINT      = "/sapi/v1/payment/fast-withdrawals/v2/initiate/"
	PROCESS_FAST_WITHDRAWAL_ENDPOINT    = "/sapi/v1/payment/fast-withdrawals/v2/process/"
	LIST_NORMAL_WITHDRAWALS_ENDPOINT    = "/sapi/v1/payment/withdrawals/"
	LIST_FAST_WITHDRAWALS_ENDPOINT      = "/sapi/v1/payment/fast-withdrawals/"

	MAINET_STARK_CONTRACT  = "0x1390f521A79BaBE99b69B37154D63D431da27A07"
	TESTNET_STARK_CONTRACT = "0xA2eC709125Ea693f5522aEfBBC3cb22fb9146B52"
)

// token/currency name
type Currency string

const (
	BTC   Currency = "btc"
	USDT  Currency = "usdt"
	ETH   Currency = "eth"
	USDC  Currency = "usdc"
	MATIC Currency = "matic"
	ARB   Currency = "arb"
	LINK  Currency = "link"
	AAVE  Currency = "aave"
	MKR   Currency = "mkr"
	LDO   Currency = "ldo"
	RPL   Currency = "rpl"
	STRK  Currency = "strk"
	BERA  Currency = "bera"
)

// network name
type Network string

const (
	ETHEREUM      Network = "ETHEREUM"
	POLYGON       Network = "POLYGON"
	STARKNET      Network = "STARKNET"
	SCROLL        Network = "SCROLL"
	OPTIMISM      Network = "OPTIMISM"
	ARBITRUM      Network = "ARBITRUM"
	LINEA         Network = "LINEA"
	MODE          Network = "MODE"
	POLYGON_ZKEVM Network = "ZKPOLY"
	BERA_NETWORK  Network = "BERA"
)

type Client struct {
	httpClient   *http.Client
	jwtToken     string
	refreshToken string

	baseURL        *url.URL
	healthURL      *url.URL
	tickerURL      *url.URL
	candleStickURL *url.URL
	orderBookURL   *url.URL
	tradesURL      *url.URL

	// login
	nonceURL        *url.URL
	loginURL        *url.URL
	refreshTokenURL *url.URL

	// user
	profileURL *url.URL
	balanceURL *url.URL
	pnlURL     *url.URL

	// order
	orderURL       *url.URL
	orderNonceURL  *url.URL
	orderCreateURL *url.URL
	orderCancelURL *url.URL
	tradesListURL  *url.URL

	// deposit
	coinURL                   *url.URL
	vaultIDURL                *url.URL
	networkConfigURL          *url.URL
	cryptoDepositStartURL     *url.URL
	crossChainDepositStartURL *url.URL
	listDepositsURL           *url.URL

	// internal transfer
	internalTransferInitiateURL *url.URL
	internalTransferProcessURL  *url.URL
	internalTransferGetURL      *url.URL
	internalTransferUserURL     *url.URL
	internalTransferListURL     *url.URL

	// withdraw
	startNormalWithdrawalURL    *url.URL
	validateNormalWithdrawalURL *url.URL
	startFastWithdrawalURL      *url.URL
	processFastWithdrawalURL    *url.URL
	listNormalWithdrawalURL     *url.URL
	listFastWithdrawalURL       *url.URL
}

// create new client
func New(base BaseURL) (*Client, error) {
	baseurl, err := url.Parse(string(base))
	if err != nil {
		return nil, err
	}

	healthurl := baseurl.JoinPath(HEALTH_ENDPOINT)
	tickerurl := baseurl.JoinPath(TICKER_ENDPOINT)
	candleStickurl := baseurl.JoinPath(CANDLESTICK_ENDPOINT)
	orderBookurl := baseurl.JoinPath(ORDERBOOK_ENDPOINT)
	tradesurl := baseurl.JoinPath(TRADES_ENDPOINT)

	nonceurl := baseurl.JoinPath(NONCE_ENDPOINT)
	loginurl := baseurl.JoinPath(LOGIN_ENDPOINT)
	refreshTokenurl := baseurl.JoinPath(REFRESH_TOKEN_ENDPOINT)

	profileurl := baseurl.JoinPath(PROFILE_ENDPOINT)
	balanceurl := baseurl.JoinPath(BALANCE_ENDPOINT)
	pnlurl := baseurl.JoinPath(PNL_ENDPOINT)

	orderURL := baseurl.JoinPath(ORDER_BASE_ENDPOINT)
	ordernonceurl := baseurl.JoinPath(ORDER_NONCE_ENDPOINT)
	ordercreateurl := baseurl.JoinPath(ORDER_CREATE_ENDPOINT)
	ordercancelurl := baseurl.JoinPath(ORDER_CANCEL_ENDPOINT)
	tradesListurl := baseurl.JoinPath(TRADES_LIST_ENDPOINT)

	coinurl := baseurl.JoinPath(COIN_ENDPOINT)
	vaultidurl := baseurl.JoinPath(VAULTID_ENDPOINT)
	networkConfigurl := baseurl.JoinPath(NETWORK_CONFIG_ENDPOINT)
	cryptoDepositStarturl := baseurl.JoinPath(CRYPTO_DEPOSIT_START_ENDPOINT)
	crossChainDepositStarturl := baseurl.JoinPath(CROSS_CHAIN_DEPOSIT_START_ENDPOINT)
	listDepositsurl := baseurl.JoinPath(LIST_DEPOSITS_ENDPOINT)

	internalTransferInitiate := baseurl.JoinPath(INTERNAL_TRANSFER_INITIATE_ENDPOINT)
	internalTransferProcess := baseurl.JoinPath(INTERNAL_TRANSFER_PROCESS_ENDPOINT)
	internalTransferGet := baseurl.JoinPath(INTERNAL_TRANSFER_GET_ENDPOINT)
	internalTransferUser := baseurl.JoinPath(INTERNAL_TRANSFER_USER_ENDPOINT)
	internalTransferList := baseurl.JoinPath(INTERNAL_TRANSFER_LIST_ENDPOINT)

	startNormalWithdrawal := baseurl.JoinPath(START_NORMAL_WITHDRAWAL_ENDPOINT)
	validateNormalWithdrawal := baseurl.JoinPath(VALIDATE_NORMAL_WITHDRAWAL_ENDPOINT)
	startFastWithdrawal := baseurl.JoinPath(START_FAST_WITHDRAWAL_ENDPOINT)
	processFastWithdrawal := baseurl.JoinPath(PROCESS_FAST_WITHDRAWAL_ENDPOINT)
	listNormalWithdrawal := baseurl.JoinPath(LIST_NORMAL_WITHDRAWALS_ENDPOINT)
	listFastWithdrawal := baseurl.JoinPath(LIST_FAST_WITHDRAWALS_ENDPOINT)

	return &Client{
		httpClient:   http.DefaultClient,
		jwtToken:     "",
		refreshToken: "",

		baseURL:        baseurl,
		healthURL:      healthurl,
		tickerURL:      tickerurl,
		candleStickURL: candleStickurl,
		orderBookURL:   orderBookurl,
		tradesURL:      tradesurl,

		nonceURL:        nonceurl,
		loginURL:        loginurl,
		refreshTokenURL: refreshTokenurl,

		profileURL: profileurl,
		balanceURL: balanceurl,
		pnlURL:     pnlurl,

		orderURL:       orderURL,
		orderNonceURL:  ordernonceurl,
		orderCreateURL: ordercreateurl,
		orderCancelURL: ordercancelurl,
		tradesListURL:  tradesListurl,

		coinURL:                   coinurl,
		vaultIDURL:                vaultidurl,
		networkConfigURL:          networkConfigurl,
		cryptoDepositStartURL:     cryptoDepositStarturl,
		crossChainDepositStartURL: crossChainDepositStarturl,
		listDepositsURL:           listDepositsurl,

		internalTransferInitiateURL: internalTransferInitiate,
		internalTransferProcessURL:  internalTransferProcess,
		internalTransferGetURL:      internalTransferGet,
		internalTransferUserURL:     internalTransferUser,
		internalTransferListURL:     internalTransferList,

		startNormalWithdrawalURL:    startNormalWithdrawal,
		validateNormalWithdrawalURL: validateNormalWithdrawal,
		startFastWithdrawalURL:      startFastWithdrawal,
		processFastWithdrawalURL:    processFastWithdrawal,
		listNormalWithdrawalURL:     listNormalWithdrawal,
		listFastWithdrawalURL:       listFastWithdrawal,
	}, nil
}

// check auth
func (c *Client) CheckAuth() error {
	if c.jwtToken == "" || c.refreshToken == "" {
		// throw error
		return ErrNotLoggedIn
	}

	return nil
}
