package omidpay

type VerifyResponse struct {
	Result     string `json:"Result"`
	Amount     int64  `json:"Amount"`
	RefNum     string `json:"RefNum"`
	HashedPan  string `json:"HashedPan"`
	ResultCode int    `json:"ResultCode"`
}
