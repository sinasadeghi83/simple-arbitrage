package nobitex

import (
	"encoding/json"
	"fmt"

	"github.com/sinasadeghi83/simple-arbitrage/exchange"
	"github.com/sinasadeghi83/simple-arbitrage/metrics"
	"resty.dev/v3"
)

type Nobitex struct {
	cl *resty.Client
}

type OrderBook struct {
	Status         string     `json:"status"`
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

func (nob Nobitex) GetOrderBook(pair string) (*exchange.OrderBook, error) {
	res, err := nob.cl.R().
		SetPathParam(
			"pair", pair,
		).
		Get("https://apiv2.nobitex.ir/v3/orderbook/{pair}")

	metrics.ResponseTime.WithLabelValues("nobitex").Observe(float64(res.Duration().Milliseconds()))
	if err != nil {
		metrics.RequestStatus.WithLabelValues("nobitex", "failed").Inc()
		return nil, err
	}

	if res.StatusCode() != 200 {
		metrics.RequestStatus.WithLabelValues("nobitex", "failed").Inc()
		return nil, fmt.Errorf("status code: %d", res.StatusCode())
	}

	metrics.RequestStatus.WithLabelValues("nobitex", "success").Inc()
	var nobOrderBook OrderBook
	if err := json.Unmarshal(res.Bytes(), &nobOrderBook); err != nil {
		return nil, err
	}

	orderBook := exchange.OrderBook{
		Pair:       pair,
		LastUpdate: nobOrderBook.LastUpdate,
		Asks:       nobOrderBook.Asks,
		Bids:       nobOrderBook.Bids,
	}

	return &orderBook, nil
}
