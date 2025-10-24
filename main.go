package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinasadeghi83/simple-arbitrage/arbitrage"
	"github.com/sinasadeghi83/simple-arbitrage/config"
	_ "github.com/sinasadeghi83/simple-arbitrage/metrics"
	"github.com/spf13/viper"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := config.Configure(ctx)
	if err != nil {
		panic(err)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(viper.GetString("bot.token"), opts...)
	if err != nil {
		panic(err)
	}

	arbChan := make(chan arbitrage.ArbDetail)
	pairs := viper.GetStringSlice("arbitrage.pairs")
	for _, pair := range pairs {
		go arbitrage.FindPeriodicly(ctx, pair, time.Second, arbChan)
	}
	go func() {
		for ctx.Err() == nil {
			arbDetail := <-arbChan
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: viper.GetString("bot.channel_id"),
				Text: fmt.Sprintf("New arbitrage opportunity found:\n Pair: %s\n Src: %s\n Dest: %s\n Amount: %s\n Ask: %s\n Bid: %s",
					arbDetail.Pair, arbDetail.Src, arbDetail.Dest, arbDetail.Amount, arbDetail.Ask, arbDetail.Bid),
			})
			fmt.Println(arbDetail)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(":"+viper.GetString("METRICS_PORT"), nil)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

}
