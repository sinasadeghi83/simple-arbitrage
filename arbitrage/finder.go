package arbitrage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sinasadeghi83/simple-arbitrage/nobitex"
	"github.com/sinasadeghi83/simple-arbitrage/tabdeal"
)

type ArbDetail struct {
	Src        string
	Dest       string
	Amount     string
	Ask        string
	Bid        string
	LastUpdate uint
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

	// fmt.Printf("NOB TAB: %s %s\n", leastNobBid, leastTabBid)
	// fmt.Printf("TAB NOB: %s %s\n", leastTabAsk, leastNobBid)

	if leastNobAsk < leastTabBid {
		return nil, &ArbDetail{
			Src:        "Nobitex",
			Dest:       "Tabdeal",
			Amount:     min(nobRes.Asks[0][1], tabRes.Bids[0][1]),
			Ask:        leastNobAsk,
			Bid:        leastTabBid,
			LastUpdate: nobRes.LastUpdate,
		}
	}

	if leastTabAsk < leastNobBid {
		return nil, &ArbDetail{
			Src:        "Tabdeal",
			Dest:       "Nobitex",
			Amount:     min(tabRes.Asks[0][1], nobRes.Bids[0][1]),
			Ask:        leastTabAsk,
			Bid:        leastNobBid,
			LastUpdate: nobRes.LastUpdate,
		}
	}

	return ErrNoArb, nil
}

func FindPeriodicly(ctx context.Context, pair string, period time.Duration, arb chan ArbDetail) {
	nob := nobitex.New()
	tab := tabdeal.New()
	lastUpdate := uint(time.Now().UnixMilli())

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

		if arbDetail.LastUpdate > lastUpdate {
			lastUpdate = arbDetail.LastUpdate
			// fmt.Println(arbDetail)
			arb <- *arbDetail
		}
	}
}
