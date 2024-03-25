package domain

type Trade struct {
	Symbol             string  `json:"s"`
	TradeID            int64   `json:"t"`
	Price              float64 `json:"p,string"`
	Quantity           float64 `json:"q,string"`
	BuyerOrderID       int64   `json:"b"`
	SellerOrderID      int64   `json:"a"`
	TradeTime          int64   `json:"T"`
	IsBuyerMarketMaker bool    `json:"m"`
	Ignore             bool    `json:"M"`
}
type TradeService interface {
	GetTrade() (Trade, error)
}
