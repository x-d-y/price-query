package main

import (
	"context"
	"eth-service-price_query/lib/consts"
	"eth-service-price_query/service/cronQuery/history"
	"flag"
	"fmt"
	"time"
)

func main() {
	source := flag.String("source", "binance", "数据来源")
	symbol := flag.String("symbol", "ETHUSDT", "币种对")
	interval := flag.String("interval", "1m", "数据间隔")
	mode := flag.String("mod", "single", "运行模式")
	flag.Parse()
	fmt.Println(*source, *symbol, *interval)
	switch *source {
	case "binance":
		history.InitBinanceHistory(*symbol, *interval)
		if *mode == "batch" || *mode == "combine" {
			if historyErr := history.Requester.FetchHistoryBatch(0, time.Now().Unix()*1000, *symbol, consts.InsertModeBatch); historyErr != nil {
				fmt.Println(historyErr)
			}
		}

		if *mode == "single" || *mode == "combine" {
			ctx, _ := context.WithCancel(context.Background())
			history.Requester.CronFetch(ctx, *symbol)
		}

	default:
		panic("unsupported source")
	}
}
