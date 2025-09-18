package apiQuery

import (
	"context"
	"encoding/json"
	"errors"
	"eth-service-price_query/lib/db/timescalDb"
	"eth-service-price_query/lib/httpRequest"
	"eth-service-price_query/lib/utils"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Binance struct {
	baseQuery
	binanceRequest
	binanceResponse
}

type binanceRequest struct{}

type binanceResponse struct {
	Price string `json:"price"`
}

type Response [][]interface{}

func newBinance(interval string) *Binance {
	res := new(Binance)
	res.url = "https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=10"
	res.switchQuery = true //todo 从配置文件读取
	res.requester = httpRequest.NewClient()
	res.name = "binance"
	res.tradingTypes = "SPOT"
	res.interval = interval
	return res
}

func (b *Binance) source() string {
	return b.name
}

func (b *Binance) fetchHistory(ctx context.Context, startTime, endTime int64, symbol string) ([]*timescalDb.PriceTicks, error) {
	targetUrl := fmt.Sprintf(b.baseQuery.url, symbol, b.interval, startTime, endTime)
	res, rErr := b.requester.ContextGet(ctx, targetUrl, nil)
	if rErr != nil {
		return nil, rErr
	}
	defer res.Body.Close()
	bRes, resErr := io.ReadAll(res.Body)
	if resErr != nil {
		return nil, resErr
	}
	var resStruct = &Response{}
	if e := json.Unmarshal(bRes, resStruct); e != nil {
		return nil, e
	}
	inserts := make([]*timescalDb.PriceTicks, 0)
	for _, tick := range *resStruct {

		price, parsePriceError := strconv.ParseFloat(tick[4].(string), 64)
		if parsePriceError != nil {
			return nil, parsePriceError
		}
		volume, parseVolumeError := strconv.ParseFloat(tick[5].(string), 64)
		if parseVolumeError != nil {
			return nil, parsePriceError
		}
		insert := &timescalDb.PriceTicks{
			Time:   time.Unix(int64(tick[6].(float64))/1000, 0),
			Symbol: symbol,
			Price:  price,
			Source: b.source(),
			Volume: volume,
			Type:   b.types(ctx),
		}
		inserts = append(inserts, insert)
	}
	return inserts, nil
}

func (b *Binance) query(ctx context.Context) ([]*timescalDb.PriceTicks, error) {
	return b.fetchHistory(ctx)
}

func (b *Binance) stop(ctx context.Context) error {
	return nil
}

func (b *Binance) start(ctx context.Context) error {
	return nil
}

func (b *Binance) types(ctx context.Context) string {
	return b.baseQuery.tradingTypes
}
