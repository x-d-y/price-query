package apiQuery

import (
	"context"
	"encoding/json"
	"errors"
	"eth-service-price_query/lib/httpRequest"
	"io"
)

type Coinpaprika struct {
	baseQuery
	coinpaprikaRequest
	coinpaprikaResponse
}

type coinpaprikaRequest struct{}

type coinpaprikaResponse struct {
	Quotes struct {
		USD struct {
			Price float64 `json:"price"`
		} `json:"USD"`
	} `json:"quotes"`
}

func newCoinpaprika() *Coinpaprika {
	return &Coinpaprika{
		baseQuery: baseQuery{
			name:        "Coinpaprika",
			url:         "https://api.coinpaprika.com/v1/tickers/eth-ethereum",
			requester:   httpRequest.NewClient(),
			switchQuery: true,
		},
	}
}

func (b *Coinpaprika) source() string {
	return b.name
}

func (b *Coinpaprika) query(ctx context.Context) (float64, error) {
	if !b.switchQuery {
		return 0.0, errors.New("query is not switchQuery")
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
	res := &coinpaprikaResponse{}
	if jErr := json.Unmarshal(bRes, res); jErr != nil {
		return 0, jErr
	}
	return res.Quotes.USD.Price, nil
}

func (b *Coinpaprika) stop(ctx context.Context) error {
	return nil
}

func (b *Coinpaprika) start(ctx context.Context) error {
	return nil
}

func (b *Coinpaprika) types(ctx context.Context) string {
	return b.baseQuery.tradingTypes
}
