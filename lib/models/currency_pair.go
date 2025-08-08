package models

type CurrencyPairRequest struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
