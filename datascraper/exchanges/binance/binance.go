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

func Klineevents(conn *clickhouse.Conn) {
	wsCombinedKlineHandler := func(kline *binance.WsKlineEvent) {
		kline := kline.Kline
		isFinal := 1
		if kline.IsFinal {
			isFinal := 1
			query := clickhouse.NewQuery(fmt.Sprintf("INSERT INTO binance_klines.%s_closed VALUES (%d,%d,%s,%s,%d,%d,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s)", kline.StartTime, kline.EndTime, kline.Symbol, kline.Interval, kline.FirstTradeID, kline.LastTradeID, kline.Open, kline.Close, kline.High, kline.Low, kline.Volume, kline.TradeNum, IsFinal, kline.QuoteVolume, kline.ActiveBuyVolume, kline.ActiveQuoteVolume))
			err := query.Exec(conn)
			if err != nil {
				log.Info(err)
			}
		}
		if kline.isFinal == true {
			isFinal = 1
		} else {
			isFinal = 0
		}
		query := clickhouse.NewQuery(fmt.Sprintf("INSERT INTO binance_klines.%s_all (%d,%d,%s,%s,%d,%d,%s    ,%s,%s,%s,%s,%d,%d,%s,%s,%s)", kline.StartTime, kline.EndTime, kline.Symbol, kline.Interval, kline.FirstTradeID, kline.LastTradeID, kline.Open, kline.Close, kline.High, kline.Low, kline.Volume, kline.TradeNum, IsFinal, kline.QuoteVolume, kline.ActiveBuyVolume, kline.ActiveQuoteVolume))
		err := query.Exec(conn)
		if err != nil {
			log.Error(err)
		}
	}

	errHandler := func(err error) {
		log.Error(err)
	}

	doneC, _, err := binance.WsCombinedKlineServe([]string{"BTCUSDT", "ETHUSDT", "DOGEUSDT"}, wsCombinedKlineHandler, errHandler)
	if err != nil {
		log.Error(err)
		return
	}

	<-doneC

}
