package apiQuery

type CoinMarketCap struct {
	baseQuery
	coinMarketCapPriceRequest
	coinMarketCapPriceResponse
}

type coinMarketCapPriceRequest struct{}

type coinMarketCapPriceResponse struct {
	Status coinMarketCapStatus
}

type coinMarketCapStatus struct {
	Timestamp string `json:"timestamp"`
}
