package lib

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/hust-tianbo/go_lib/log"
)

var bot *openwechat.Bot

func Init() {
	bot = openwechat.DefaultBot()

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	if err := bot.Login(); err != nil {
		panic(err)
	}

	// 注释阻塞流程
	//bot.Block()
}

func ReportToWX(recv string, content string) {
	// 强制写成测试账号
	recv = "田博测试"

	if bot == nil {
		Init()
	}
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Errorf("[ReportToWX]get user failed:%+v", err)
		return
	}
	friends, err := self.Friends()
	if err != nil {
		log.Errorf("[ReportToWX]get friend failed:%+v", err)
		return
	}

	/*log.Debugf("[ReportToWX]friends size %+v", len(friends))
	log.Debugf("[ReportToWX]all friends %+v", friends)*/

	friend := friends.SearchByRemarkName(1, recv)
	if friend.Count() > 0 {
		friend.SendText(content)
	}
	log.Debugf("[ReportToWX]finish send text:%+v|%+v", recv, content)
}
