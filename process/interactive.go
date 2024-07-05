package process

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InteractFunc func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, data *interface{}) bool

type InteractiveInfo struct {
	CurState          int
	Data              *interface{} // 例如:AddRestaurantState
	InteractFuncSlice []*InteractFunc
}

var InteractList = make(map[int64]*InteractiveInfo)

// regist interactive userid
func RegisInteractiveMode(userId int64, data *interface{}, interactFunc []*InteractFunc) {
	InteractList[userId] = &InteractiveInfo{
		CurState:          0,
		Data:              data,
		InteractFuncSlice: interactFunc,
	}
}

// process Interactive
func ProcessInteractive(userId int64, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	II := InteractList[userId]
	IFS := II.InteractFuncSlice
	f := *IFS[II.CurState]
	if f(bot, update, II.Data) {
		II.CurState++
	}
	if II.CurState > len(II.InteractFuncSlice) {
		DelInteractive(userId)
	}
}

func DelInteractive(userId int64) {
	delete(InteractList, userId)
}
