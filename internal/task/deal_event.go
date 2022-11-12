package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hust-tianbo/go_lib/log"
	"github.com/hust-tianbo/stock_monitor/internal/logic"
	"github.com/hust-tianbo/stock_monitor/internal/model"
	"github.com/hust-tianbo/stock_monitor/lib"
)

func UpdateEventStatus(event *model.EventList, status int) {
	now := time.Now()
	logic.EventDb.Table(model.EventTable).Where(&model.EventList{EventName: event.EventName}).Update(map[string]interface{}{"status": status, "m_time": now})
}

func GetDoingEventList() ([]model.EventList, error) {
	var eventList []model.EventList
	dbRes := logic.EventDb.Table(model.EventTable).Where("status=?", model.EventStatusInit).Find(&eventList)

	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[GetDoingEventList]find record failed:%+v", dbRes.Error)
		return eventList, dbRes.Error
	}

	return eventList, nil
}

// 示例：1668269633
func DealTimeEvent(event *model.EventList) {
	triggerTime, parseError := strconv.ParseInt(event.Extra, 10, 64)
	if parseError != nil {
		log.Errorf("[DealTimeEvent]parse failed:%+v|%+v", event, parseError)
		return
	}

	now := time.Now().Unix()

	if now >= triggerTime {
		lib.ReportToWX(event.RecvUser, fmt.Sprintf("%+v", event.EventName))
	}

	UpdateEventStatus(event, model.EventStatusFinish)

	return
}

// 示例：00700_250.0_280.0
func DealStockEvent(event *model.EventList) {
	stockParam := strings.Split(event.Extra, "_")
	if len(stockParam) != 3 {
		log.Errorf("[DealStockEvent]extra length invalid:%+v", event)
		return
	}

	stockName := stockParam[0]
	minPrice, minError := strconv.ParseFloat(stockParam[1], 32)
	maxPrice, maxError := strconv.ParseFloat(stockParam[2], 32)
	curPrice, getError := lib.GetStockPrice(stockName)
	if minError != nil || maxError != nil || getError != nil {
		log.Errorf("[DealStockEvent]parse failed:%+v|%+v|%+v", minError, maxError, getError)
		return
	}

	if curPrice <= float32(minPrice) {
		lib.ReportToWX(event.RecvUser, fmt.Sprintf(
			"告警项:%+v 最小价格:%+v 当前价格:%+v", event.EventName, minPrice, curPrice))
	}

	if curPrice >= float32(maxPrice) {
		lib.ReportToWX(event.RecvUser, fmt.Sprintf(
			"告警项:%+v 最大价格:%+v 当前价格:%+v", event.EventName, maxPrice, curPrice))
	}

	UpdateEventStatus(event, model.EventStatusFinish)

	return
}

func HandleEvent() {
	list, err := GetDoingEventList()
	if err != nil {
		return
	}

	for _, ele := range list {
		switch ele.EventType {
		case logic.EventType_Stock:
			DealStockEvent(&ele)
			break
		case logic.EventType_Time:
			DealTimeEvent(&ele)
			break
		default:
			log.Errorf("[HandleEvent]event type invalid:%+v", ele)
		}
	}
}

func DoingTask() {
	ticket := time.NewTicker(time.Minute * 3)
	go func() {
		for {
			select {
			case <-ticket.C:
				HandleEvent()
			}
		}
	}()
}
