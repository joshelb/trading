package clickhouse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func get_symbols() []string {
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var symbols map[string]interface{}

	err = json.Unmarshal(body, &symbols)
	if err != nil {
		log.Error(err)
	}

	var result []string
	//log.Info(symbols["symbols"].([]interface{})[0].(map[string]interface{})["symbol"])
	for _, pair := range symbols["symbols"].([]interface{}) {
		result = append(result, strings.ToLower(pair.(map[string]interface{})["symbol"].(string)))
	}
	return result
}
