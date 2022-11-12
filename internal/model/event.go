package model

import "time"

const EventTable string = "event_list"

const (
	EventStatusInit   = 0
	EventStatusFinish = 1
	EventStatusDel    = 2
)

type EventList struct {
	EventName string    `gorm:"column:event_name"`
	EventType int       `gorm:"column:event_type"`
	Extra     string    `gorm:"column:extra"`
	OpeUser   string    `gorm:"column:ope_user"`
	RecvUser  string    `gorm:"column:recv_user"`
	CTime     time.Time `gorm:"column:c_time"`
	MTime     time.Time `gorm:"column:m_time"`
	Status    int       `gorm:"column:status"` // 0.新建 1.已经触发过 2.已删除
}
