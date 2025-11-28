package entity

import "time"

type Trade struct {
	ID        string    `json:"id" bson:"_id"`
	Symbol    string    `json:"symbol" bson:"symbol"`
	Price     float64   `json:"price" bson:"price"`
	Quantity  float64   `json:"quantity" bson:"quantity"`
	TradeTime time.Time `json:"trade_time" bson:"trade_time"`
	IsBuyer   bool      `json:"is_buyer" bson:"is_buyer"`
}

type BinanceAggTrade struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int64  `json:"a"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int64  `json:"T"`
	IsBuyer   bool   `json:"m"`
}
