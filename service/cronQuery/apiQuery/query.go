package apiQuery

import (
	"context"
	"eth-service-price_query/lib/consts"
	"eth-service-price_query/lib/db/timescalDb"
	"eth-service-price_query/lib/httpRequest"
	"fmt"
	"time"
)

func Init() {
	for {
		now := time.Now() // 多个查询的时候，使用同一个时间戳， 方便后面数据汇总
		processQuery(now) // todo 改成 goroutine 查多个
		time.Sleep(10 * time.Second)
	}
}

func processQuery(now time.Time) {
	for _, v := range apiList {
		ctx := context.Background()
		price, priceErr := v.query(ctx)
		if priceErr != nil {
			//todo log 记录
			fmt.Println(priceErr)
			continue
		}
		pt := &timescalDb.PriceTicks{
			Time:   now,
			Symbol: consts.ETH,
			Price:  price,
			Source: v.source(),
		}
		fmt.Println(pt) //todo delete me
		if dbErr := timescalDb.TimescaleDb.InsertPriceTick(ctx, pt); dbErr != nil {
			fmt.Println(dbErr) // todo 日志记录
		}
	}
}

var apiList = []Price{
	newBinance(),     // 币安查询
	newCoinpaprika(), // coinPaprika 查询
}

type Price interface {
	query(ctx context.Context) (float64, error)
	source() string
	stop(ctx context.Context) error
	start(ctx context.Context) error
	types(ctx context.Context) string
}

type baseQuery struct {
	name         string
	url          string
	switchQuery  bool
	requester    *httpRequest.Requester
	tradingTypes string
	interval     string
}
