package models

type UniswapTransaction struct {
	From      string `json:"from"`
	To        string `json:"to"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	Timestamp int64  `json:"timestamp"`
}
