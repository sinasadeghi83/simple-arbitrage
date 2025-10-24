package arbitrage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sinasadeghi83/simple-arbitrage/metrics"
	"github.com/sinasadeghi83/simple-arbitrage/nobitex"
	"github.com/sinasadeghi83/simple-arbitrage/tabdeal"
)

type ArbDetail struct {
	Pair   string
	Src    string
	Dest   string
	Amount string
	Ask    string
	Bid    string
}

var ErrNoArb = fmt.Errorf("No arbitrage found")

func Find(pair string, nob nobitex.Nobitex, tab tabdeal.Tabdeal) (error, *ArbDetail) {
	nobRes, err := nob.GetOrderBook(pair)
	if err != nil {
		return err, nil
	}

	tabRes, err := tab.GetOrderBook(pair)
	if err != nil {
		return err, nil
	}

	leastNobBid := nobRes.Bids[0][0]
	leastTabBid := tabRes.Bids[0][0]

	leastNobAsk := nobRes.Asks[0][0]
	leastTabAsk := tabRes.Asks[0][0]

	if leastNobAsk < leastTabBid {
		metrics.ArbitrageOpportunityRate.Inc()
		askPrice, _ := strconv.ParseFloat(leastNobAsk, 64)
		bidPrice, _ := strconv.ParseFloat(leastTabBid, 64)
		metrics.LastPriceDifference.WithLabelValues(pair, "Nobitex", "Tabdeal").Set(bidPrice - askPrice)

		return nil, &ArbDetail{
			Pair:   pair,
			Src:    "Nobitex",
			Dest:   "Tabdeal",
			Amount: min(nobRes.Asks[0][1], tabRes.Bids[0][1]),
			Ask:    leastNobAsk,
			Bid:    leastTabBid,
		}
	}

	if leastTabAsk < leastNobBid {
		metrics.ArbitrageOpportunityRate.Inc()
		askPrice, _ := strconv.ParseFloat(leastTabAsk, 64)
		bidPrice, _ := strconv.ParseFloat(leastNobBid, 64)
		metrics.LastPriceDifference.WithLabelValues(pair, "Tabdeal", "Nobitex").Set(bidPrice - askPrice)

		return nil, &ArbDetail{
			Pair:   pair,
			Src:    "Tabdeal",
			Dest:   "Nobitex",
			Amount: min(tabRes.Asks[0][1], nobRes.Bids[0][1]),
			Ask:    leastTabAsk,
			Bid:    leastNobBid,
		}
	}

	return ErrNoArb, nil
}

func FindPeriodicly(ctx context.Context, pair string, period time.Duration, arb chan ArbDetail) {
	nob := nobitex.New()
	tab := tabdeal.New()
	lastArb := ArbDetail{}

	for ctx.Err() == nil {
		time.Sleep(time.Second)
		err, arbDetail := Find(pair, nob, tab)
		if err != nil {
			if !errors.Is(err, ErrNoArb) {
				close(arb)
				fmt.Println(err)
				return
			}
			continue
		}

		if *arbDetail != lastArb {
			arb <- *arbDetail
			lastArb = *arbDetail
		}
	}
}
