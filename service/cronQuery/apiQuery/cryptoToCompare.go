package apiQuery

type cryptoToCompare struct {
	baseQuery
	cryptoToCompareRequest
	cryptoToCompareResponse
}

type cryptoToCompareRequest struct{}

type cryptoToCompareResponse struct {
	USD float64 `json:"USD"`
}
