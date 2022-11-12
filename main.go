package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/hust-tianbo/stock_monitor/config"
	"github.com/hust-tianbo/stock_monitor/internal/logic"
)

const (
	port = ":51050"
)

func main() {
	//modbus.get()
	log.Debugf("begin logic server")
	config.InitConfig()

	logic.InitImp()

	// 注册http接口
	mux := GetHttpServerMux()
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}),
	)

	log.Debugf("end logic server")
}

func GetHttpServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/get_box_info", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.AddEventReq
		json.Unmarshal(body, &req)
		var rsp logic.AddEventRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = *logic.AddEvent(&req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	return mux
}
