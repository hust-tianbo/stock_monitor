package logic

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hust-tianbo/stock_monitor/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var EventDb *gorm.DB

func InitImp() {
	cf := config.GetConfig()
	// 随机数种子
	rand.Seed(time.Now().Unix())

	// 连接db
	var err error
	EventDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/register_event?charset=utf8&parseTime=True&loc=Local", cf.DBUser, cf.DBSecret, cf.DBIP))
	if err != nil {
		panic(err)
	}
	EventDb.SingularTable(true)
}
