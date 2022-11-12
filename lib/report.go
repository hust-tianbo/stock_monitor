package lib

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/hust-tianbo/go_lib/log"
)

var bot *openwechat.Bot

func init() {
	Init()
}

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

	bot.Block()
}

func ReportToWX(recv string, content string) {
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

	friend := friends.SearchByUserName(1, recv)
	if friend.Count() > 0 {
		friend.SendText(content)
	}

}
