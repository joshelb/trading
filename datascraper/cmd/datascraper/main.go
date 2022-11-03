package main

import (
	"fmt"

	"github.com/joshelb/datascraper/clickhouse"
	"github.com/joshelb/datascraper/exchanges/binance"
)

func main() {
	conn, err := clickhouse.Connect()
	if err != nil {
		fmt.Print(err)
	}
	binance.Tradeevents(conn)
}
