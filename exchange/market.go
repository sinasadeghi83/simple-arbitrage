package exchange

type OrderBook struct {
	Pair       string     `json:"pair"`
	LastUpdate uint       `json:"lastUpdate"`
	Asks       [][]string `json:"asks"`
	Bids       [][]string `json:"bids"`
}
