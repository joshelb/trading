package binance

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/roistat/go-clickhouse"
	log "github.com/sirupsen/logrus"
)

func Tradeevents(conn *clickhouse.Conn) {
	wsCombinedTradeHandler := func(trade *binance.WsCombinedTradeEvent) {
		event := trade.Data
		stream := trade.Stream[:len(trade.Stream)-6]
		log.Info(stream)
		var IsBuyerMaker int
		if event.IsBuyerMaker == true {
			IsBuyerMaker = 1
		} else {
			IsBuyerMaker = 0
		}
		query, err := clickhouse.BuildInsert(fmt.Sprintf("binance_trades.%sbuffer", stream),
			clickhouse.Columns{"eventtype", "eventtime", "symbol", "tradeID", "price", "quantity", "buyerOrderID", "sellerOrderID", "tradeTime", "isBuyerMaker", "placeholder", "time"},
			clickhouse.Row{event.Event, event.Time, event.Symbol, event.TradeID, event.Price, event.Quantity, event.BuyerOrderID, event.SellerOrderID, event.TradeTime, IsBuyerMaker, 1, (time.Now().String())},
		)
		if err != nil {
			log.Error(err)
		}
		err = query.Exec(conn)
		if err != nil {
			log.Error(err)
		}
	}

	errHandler := func(err error) {
		log.Error(err)
	}

	doneC, _, err := binance.WsCombinedTradeServe([]string{"BTCUSDT", "ETHUSDT", "DOGEUSDT"}, wsCombinedTradeHandler, errHandler)
	if err != nil {
		log.Error(err)
		return
	}

	<-doneC

}
