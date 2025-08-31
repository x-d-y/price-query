package history

import (
	"context"
	"encoding/json"
	"eth-service-price_query/lib/db/timescalDb"
	"eth-service-price_query/lib/httpRequest"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type binanceRequester struct {
	*httpRequest.Requester
	targetUrl string
	source    string
	types     string
	interval  string
	symbol    string
}

var Requester *binanceRequester

func InitBinanceHistory(symbol, interval, types string) {
	Requester = &binanceRequester{
		Requester: httpRequest.NewClient(),
		targetUrl: "https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=1000",
		source:    "binance",
		types:     types,
		interval:  interval,
		symbol:    symbol,
	}
}

//[
//1499040000000,      // Open time (ms)
//"0.01634790",       // Open
//"0.80000000",       // High
//"0.01575800",       // Low
//"0.01577100",       // Close
//"148976.11427815",  // Volume
//1499644799999,      // Close time (ms)
//"2434.19055334",    // Quote asset volume
//308,                // Number of trades
//"1756.87402397",    // Taker buy base asset volume
//"28.46694368",      // Taker buy quote asset volume
//"17928899.62484339" // Ignore (can ignore)
//]

type Response [][]interface{}

func (r *binanceRequester) FetchHistoryBatch(startTime, endTime int64, symbol string) error {
	for endTime > startTime {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		inserts, insertErr := r.fetchHistory(ctx, startTime, endTime, symbol)
		if insertErr != nil {
			return insertErr
		}

		if dbErr := timescalDb.TimescaleDb.InsertPriceTickBatch(ctx, inserts, r.interval); dbErr != nil {
			if !strings.Contains(dbErr.Error(), "SQLSTATE 23505") {
				return dbErr
			}
		}

		if len(inserts) < 1000 { // 当数量不够的时候说明没有更多的数据了，需要退出
			return nil
		}

		startTime = inserts[len(inserts)-1].Time.Unix()*1000 + 1

		if startTime >= endTime { // 左闭右开， 时间超出， 结束
			return nil
		}
		time.Sleep(time.Millisecond * 500)

	}
	return nil
}

func (r *binanceRequester) fetchHistory(ctx context.Context, startTime, endTime int64, symbol string) ([]*timescalDb.PriceTicks, error) {
	targetUrl := fmt.Sprintf(r.targetUrl, symbol, r.interval, startTime, endTime)
	res, rErr := r.Requester.ContextGet(ctx, targetUrl, nil)
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
			Source: r.source,
			Volume: volume,
			Type:   r.types,
		}
		inserts = append(inserts, insert)
	}
	fmt.Println(inserts)
	return inserts, nil
}
