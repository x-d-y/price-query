package main

import (
	"eth-service-price_query/service/cronQuery/history"
	"flag"
	"fmt"
	"time"
)

func main() {
	source := flag.String("source", "binance", "数据来源")
	symbol := flag.String("symbol", "ETHUSDT", "币种对")
	types := flag.String("type", "spot", "现货/期货")
	interval := flag.String("interval", "1m", "数据间隔")
	flag.Parse()
	fmt.Println(*source, *symbol, *types, *interval)
	switch *source {
	case "binance":
		history.InitBinanceHistory(*symbol, *interval, *types)
		if historyErr := history.Requester.FetchHistoryBatch(0, time.Now().Unix()*1000, *symbol); historyErr != nil {
			fmt.Println(historyErr)
		}
	default:
		panic("unsupported source")
	}
}
