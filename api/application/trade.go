// application/trade.go
package application

import (
	"log"
	"sort"

	"github.com/moasadi/binance-trade/api/domain"
)

type TradeApp struct {
	service domain.TradeService
	prices  []float64
}

func NewTradeApp(service domain.TradeService) *TradeApp {
	return &TradeApp{
		service: service,
		prices:  []float64{},
	}
}

func (app *TradeApp) Run(medianPriceChan chan<- float64) error {
	for {
		trade, err := app.service.GetTrade()
		if err != nil {
			return err
		}

		app.prices = append(app.prices, trade.Price)

		if len(app.prices) > 0 {
			sort.Float64s(app.prices)
			medianPrice := app.prices[(len(app.prices)-1)/2]

			medianPriceChan <- medianPrice
		} else {
			log.Println("No prices available to calculate median")
		}
	}

}
