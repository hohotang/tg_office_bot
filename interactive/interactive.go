package interactive

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AskFunc func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, data *interface{}) bool
type ReplyFunc func(bot *tgbotapi.BotAPI, update *tgbotapi.Update, data *interface{}) bool

type InteractiveStep struct {
	Ask   AskFunc
	Reply ReplyFunc
}

type InteractiveInfo struct {
	CurState      int
	Data          interface{} // 例如:AddRestaurantState
	InteractSteps []*InteractiveStep
}

var InteractList = make(map[int64]*InteractiveInfo)

// regist interactive userid
func RegisInteractiveMode(userId int64, data interface{}, steps []*InteractiveStep) {
	InteractList[userId] = &InteractiveInfo{
		CurState:      0,
		Data:          data,
		InteractSteps: steps,
	}
}

// process Interactive
func ProcessInteractive(userId int64, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	II := InteractList[userId]
	if II == nil {
		return
	}

	currentStep := II.InteractSteps[II.CurState]
	if currentStep.Reply != nil && currentStep.Reply(bot, update, &II.Data) {
		II.CurState++
		if II.CurState < len(II.InteractSteps) {
			nextStep := II.InteractSteps[II.CurState]
			if nextStep.Ask != nil {
				nextStep.Ask(bot, update, &II.Data)
			}
		}
	}

	if II.CurState >= len(II.InteractSteps) {
		DelInteractive(userId)
	}
}

func DelInteractive(userId int64) {
	delete(InteractList, userId)
}

func IsInInteractiveMode(userId int64) bool {
	_, ok := InteractList[userId]
	return ok
}
