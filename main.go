package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sinasadeghi83/simple-arbitrage/arbitrage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arbChan := make(chan arbitrage.ArbDetail)
	go arbitrage.FindPeriodicly(ctx, "TONIRT", time.Second, arbChan)
	go func() {
		for ctx.Err() == nil {
			arbDetail := <-arbChan
			fmt.Println(arbDetail)
		}
	}()
	<-ctx.Done()

	fmt.Print("Done")
}
