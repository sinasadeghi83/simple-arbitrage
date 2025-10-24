package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "exchange_response_time_seconds",
		Help: "Response time to the exchange for each request.",
	}, []string{"exchange"})

	RequestStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "exchange_requests_total",
		Help: "Number of successful and failed requests to each exchange.",
	}, []string{"exchange", "status"})

	ArbitrageOpportunityRate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "arbitrage_opportunities_total",
		Help: "The total number of arbitrage opportunities found.",
	})

	LastPriceDifference = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "last_price_difference",
		Help: "Last price difference for each pair.",
	}, []string{"pair", "src", "dest"})
)
