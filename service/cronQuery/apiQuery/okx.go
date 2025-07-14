package apiQuery

type Okx struct {
	baseQuery
	okxRequest
	okxResponse
}

type okxRequest struct{}

type okxResponse struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data OkxData `json:"data"`
}

type OkxData struct {
	Last float64 `json:"last"`
}
