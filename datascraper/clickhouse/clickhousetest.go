package clickhouse

import (
	"fmt"

	"github.com/roistat/go-clickhouse"
	log "github.com/sirupsen/logrus"
)

func init_tables(conn *clickhouse.Conn) {
	symbols := get_symbols()
	log.Info(symbols)
	for _, pair := range symbols {

		qdb := clickhouse.NewQuery(fmt.Sprintf("CREATE TABLE binance_trades.%s(eventtype String,eventtime Int64,symbol String,tradeID Int64,price String,quantity String,buyerOrderID Int64,sellerOrderID Int64,tradeTime Int64,isBuyerMaker Bool,placeholder Bool,time String) ENGINE = MergeTree() PRIMARY KEY (symbol, eventtime)", pair))
		qdbbuffer := clickhouse.NewQuery(fmt.Sprintf("CREATE TABLE binance_trades.%sbuffer(eventtype String,eventtime Int64,symbol String,tradeID Int64,price String,quantity String,buyerOrderID Int64,sellerOrderID Int64,tradeTime Int64,isBuyerMaker Bool,placeholder Bool,time String) ENGINE = Buffer('binance_trades','%s',16, 5, 30, 1000, 10000, 1000000, 10000000)", pair))
		err := qdb.Exec(conn)
		if err != nil {
			log.Error(err)
		}
		err = qdbbuffer.Exec(conn)
		if err != nil {
			log.Error(err)
		}
	}
}

func Connect() (*clickhouse.Conn, error) {
	transport := clickhouse.NewHttpTransport()
	conn := clickhouse.NewConn("localhost:8123", transport)
	err := conn.Ping()
	if err != nil {
		panic(err)
	}
	init_tables(conn)
	return conn, nil
}
