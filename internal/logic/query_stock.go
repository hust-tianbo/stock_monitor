package logic

import (
	"github.com/hust-tianbo/go_lib/log"
	"github.com/hust-tianbo/stock_monitor/lib"
)

type QueryStockReq struct {
	StockNum string `json:"stock_num"`
	SendUser string `json:"send_user"`
}

type QueryStockRsp struct {
	Ret   int     `json:"ret"`   // 错误码
	Msg   string  `json:"msg"`   // 错误信息
	Price float32 `json:"price"` // 当前价格
}

func QueryStock(req *QueryStockReq) *QueryStockRsp {
	var rsp = &QueryStockRsp{Ret: lib.RetSuccess}

	price, err := lib.GetStockPrice(req.StockNum)
	if err != nil {
		log.Errorf("[QueryStock]GetStockPrice failed:%+v|%+v", req, err)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	rsp.Price = price

	lib.ReportToWX(req.SendUser, "hello")
	return rsp
}
