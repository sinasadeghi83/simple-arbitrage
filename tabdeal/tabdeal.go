package tabdeal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sinasadeghi83/simple-arbitrage/exchange"
	"resty.dev/v3"
)

type Tabdeal struct {
	cl *resty.Client
}

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

func New() Tabdeal {
	cl := resty.New()
	return Tabdeal{
		cl: cl,
	}
}

func (tab Tabdeal) GetOrderBook(pair string) (*exchange.OrderBook, error) {
	res, err := tab.cl.R().
		SetQueryParam(
			"symbol", pair,
		).
		Get("https://api1.tabdeal.org/r/api/v1/depth")

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode())
	}

	var tabOrderBook OrderBook
	if err := json.Unmarshal(res.Bytes(), &tabOrderBook); err != nil {
		return nil, err
	}

	orderBook := exchange.OrderBook{
		Pair:       pair,
		LastUpdate: uint(time.Now().UnixMilli()),
		Asks:       tabOrderBook.Asks,
		Bids:       tabOrderBook.Bids,
	}

	return &orderBook, nil
}
