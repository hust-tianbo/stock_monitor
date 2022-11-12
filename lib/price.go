package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hust-tianbo/go_lib/log"
)

type StockResult struct {
	Code     string `json:"code"`
	Type     string `json:"type"`
	Market   string `json:"market"`
	Price    string `json:"price"`
	Increase string `json:"increase"`
	Ratio    string `json:"ratio"`
}

type ResultContent struct {
	Stock []StockResult `json:"stock"`
}

type FindResult struct {
	QueryID    string        `json:"QueryID"`
	ResultCode string        `json:"ResultCode"`
	Result     ResultContent `json:"result"`
}

func GetStockPrice(stockNum string) (float32, error) {
	url := fmt.Sprintf("https://finance.pae.baidu.com/selfselect/sug?wd=%s&skip_login=1&finClientType=pc", stockNum)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("[GetStockPrice]req failed:%+v|%+v", stockNum, err)
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[GetStockPrice]read content failed:%+v|%+v", stockNum, err)
		return 0, err
	}

	var res FindResult
	json.Unmarshal(body, &res)

	log.Debugf("[GetStockPrice]get result success:%+v|%+v", stockNum, res)

	if len(res.Result.Stock) < 1 {
		return 0, fmt.Errorf("not enough result")
	}

	if res.Result.Stock[0].Code != stockNum {
		return 0, fmt.Errorf("code not find")
	}
	priceFloat, err := strconv.ParseFloat(res.Result.Stock[0].Price, 32)
	if err != nil {
		return 0, err
	}
	return float32(priceFloat), nil
}
