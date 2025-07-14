package apiQuery

import (
	"context"
	"encoding/json"
	"errors"
	"eth-service-price_query/lib/httpRequest"
	"eth-service-price_query/lib/utils"
	"io"
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

func newBinance() *Binance {
	res := new(Binance)
	res.url = "https://api.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
	res.switchQuery = true //todo 从配置文件读取
	res.requester = httpRequest.NewClient()
	res.name = "binance"
	return res
}

func (b *Binance) source() string {
	return b.name
}

func (b *Binance) query(ctx context.Context) (float64, error) {
	if !b.switchQuery {
		return 0, errors.New("[binance] switch query not on")
	}

	r, rErr := b.requester.ContextGet(ctx, b.url, nil)
	if rErr != nil {
		return 0, rErr
	}
	defer r.Body.Close()

	bRes, bErr := io.ReadAll(r.Body)
	if bErr != nil {
		return 0, bErr
	}

	res := binanceResponse{}
	if jErr := json.Unmarshal(bRes, &res); jErr != nil {
		return 0, jErr
	}

	return utils.String2Float64(res.Price)
}

func (b *Binance) stop(ctx context.Context) error {
	return nil
}

func (b *Binance) start(ctx context.Context) error {
	return nil
}
