package infrastructure

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/moasadi/binance-trade/api/domain"
)

type TradeService struct {
	conn *websocket.Conn
}

func NewTradeService(conn *websocket.Conn) *TradeService {
	return &TradeService{
		conn: conn,
	}
}

func (s *TradeService) GetTrade() (domain.Trade, error) {
	_, message, err := s.conn.ReadMessage()
	if err != nil {
		return domain.Trade{}, err
	}

	var trade domain.Trade
	err = json.Unmarshal(message, &trade)
	if err != nil {
		return domain.Trade{}, err
	}

	return trade, nil
}
