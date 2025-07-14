package apiQuery

type Coinpaprika struct {
	baseQuery
	coinpaprikaRequest
	coinpaprikaResponse
}

type coinpaprikaRequest struct{}

type coinpaprikaResponse struct {
	Quotes struct {
		Price float64 `json:"price"`
	}
}
