package logic

import (
	"time"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/hust-tianbo/stock_monitor/internal/model"
	"github.com/hust-tianbo/stock_monitor/lib"
)

const (
	EventType_Stock = 1 // 股票低于价格推送
	EventType_Time  = 2 // 定时推送
)

type AddEventReq struct {
	EventName string `json:"event_name"`
	OpeUser   string `json:"ope_user"`
	EventType int    `json:"event_type"`
	RecvUser  string `json:"recv_user"`
	Extra     string `json:"extra"`
}

type AddEventRsp struct {
	Ret int    `json:"ret"` // 错误码
	Msg string `json:"msg"` // 错误信息
}

func AddEvent(req *AddEventReq) *AddEventRsp {
	var rsp = &AddEventRsp{Ret: 0}

	now := time.Now()
	dbRes := EventDb.Table(model.EventTable).Create(&model.EventList{
		EventName: req.EventName,
		EventType: req.EventType,
		Extra:     req.Extra,
		OpeUser:   req.OpeUser,
		RecvUser:  req.RecvUser,
		CTime:     now,
		MTime:     now,
		Status:    model.EventStatusInit,
	})

	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[AddEvent]create event failed:%+v,%+v", req, dbRes.Error)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	log.Debugf("[AddEvent]create event success:%+v", req)
	return rsp
}

type DelEventReq struct {
	EventName string `json:"event_name"`
}

type DelEventRsp struct {
	Ret int    `json:"ret"` // 错误码
	Msg string `json:"msg"` // 错误信息
}

func DelEvent(req *DelEventReq) *DelEventRsp {
	var rsp = &DelEventRsp{Ret: lib.RetSuccess}

	// 状态全部更新为删除状态
	now := time.Now()
	dbRes := EventDb.Table(model.EventTable).Where(&model.EventList{EventName: req.EventName}).Update(map[string]interface{}{"status": model.EventStatusDel, "m_time": now})

	if dbRes.Error != nil {
		log.Errorf("[DelEvent]del event:%+v|%+v", req, dbRes.Error)
		rsp.Ret = lib.RetInternalError
		return rsp
	}

	return rsp
}
