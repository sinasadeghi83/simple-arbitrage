package nobitex

import (
	"encoding/json"
	"fmt"

	"resty.dev/v3"
)

type Nobitex struct {
	cl *resty.Client
}

type OrderBook struct {
	Status         string     `json:"status"`
	Pair           string     `json:"pair"`
	LastUpdate     uint       `json:"lastUpdate"`
	LastTradePrice string     `json:"lastTradePrice"`
	Asks           [][]string `json:"asks"`
	Bids           [][]string `json:"bids"`
}

func New() Nobitex {
	cl := resty.New()
	return Nobitex{
		cl: cl,
	}
}

func (nob Nobitex) GetOrderBook(pair string) (*OrderBook, error) {
	res, err := nob.cl.R().
		SetPathParam(
			"pair", pair,
		).
		Get("https://apiv2.nobitex.ir/v3/orderbook/{pair}")

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode())
	}

	var orderBook OrderBook
	if err := json.Unmarshal(res.Bytes(), &orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}
