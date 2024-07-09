package process

import (
	"fmt"
	"strconv"
	"strings"
	"tgbot/constant"
	"tgbot/data"
	"tgbot/interactive"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QAInfo struct {
	Question string
	Answer   string
}

// interactive mode to put qa
func askQuestion(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	msg := tgbotapi.NewMessage(utils.GetFromID(update), "請輸入問題:")
	utils.Send(bot, msg)
	return true
}

func replyQuestion(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	info := (*qaDataPtr).(*QAInfo)
	if update.Message == nil {
		AddingFailed(bot, utils.GetFromID(update), "錯誤")
		return false
	}
	info.Question = update.Message.Text
	return true
}

func askAnswer(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	msg := tgbotapi.NewMessage(utils.GetFromID(update), "請輸入答案:")
	utils.Send(bot, msg)
	return true
}

func replyAnswer(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	info := (*qaDataPtr).(*QAInfo)
	if update.Message == nil {
		AddingFailed(bot, utils.GetFromID(update), "錯誤")
		return false
	}
	info.Answer = update.Message.Text
	return true
}

func askQAConfirm(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	info, ok := (*qaDataPtr).(*QAInfo)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析狀態")
		return false
	}
	var builder strings.Builder
	builder.WriteString("請確認以下餐廳信息:")
	builder.WriteString(fmt.Sprintf("\n\n問題: %s", info.Question))
	builder.WriteString(fmt.Sprintf("\n回答: %s", info.Answer))
	// 確認餐廳信息
	msg := tgbotapi.NewMessage(utils.GetFromID(update), builder.String())
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("確認", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_CONFIRM)),
			tgbotapi.NewInlineKeyboardButtonData("取消", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_DENY)),
		),
	)
	utils.Send(bot, msg)
	return true
}

func replyQAConfirm(bot *tgbotapi.BotAPI, update *tgbotapi.Update, qaDataPtr *interface{}) bool {
	info, ok := (*qaDataPtr).(*QAInfo)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析狀態")
		return false
	}

	if update.CallbackQuery == nil {
		AddingFailed(bot, utils.GetFromID(update), "請使用選項")
	}
	callbackData := update.CallbackQuery.Data
	str := utils.AnalysisAddRestaurantCB(callbackData)
	if str == constant.CALLBACK_CONFIRM {
		AddingQASuccess(bot, update, info)
	} else {
		AddingFailed(bot, utils.GetFromID(update), "新增已取消")
	}

	return true
}

func regisQA(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	data := &QAInfo{}
	steps := []*interactive.InteractiveStep{
		{Ask: askQuestion, Reply: replyQuestion},
		{Ask: askAnswer, Reply: replyAnswer},
		{Ask: askQAConfirm, Reply: replyQAConfirm},
	}
	interactive.RegisInteractiveMode(update.Message.From.ID, data, steps)
	steps[0].Ask(bot, update, nil)
}

func AddingQASuccess(bot *tgbotapi.BotAPI, update *tgbotapi.Update, info *QAInfo) {
	cb := update.CallbackQuery
	userName := fmt.Sprintf("%s %s", cb.From.FirstName, cb.From.LastName)
	userId := cb.From.ID
	dataQAInfo := data.QAInfo{}
	dataQAInfo.Question = info.Question
	dataQAInfo.Answer = info.Answer

	data.QAList = append(data.QAList, dataQAInfo)
	sendMessageAndClear(bot, userId, fmt.Sprintf("感謝 %s, 成功加入QA", userName))
}

func showQA(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var buttons []tgbotapi.InlineKeyboardButton
	for i, qa := range data.QAList {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", i+1, qa.Question), fmt.Sprintf("%s_%d", constant.CALLBACK_QA, i)))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "請選擇你想查看的問題:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	utils.Send(bot, msg)
}

func handleQAButtonClick(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery
	callbackData := cb.Data
	parts := strings.Split(callbackData, "_")
	if len(parts) != 2 {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的指令"))
		return
	}

	command, indexStr := parts[0], parts[1]
	if command != constant.CALLBACK_QA {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的操作"))
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的索引"))
		return
	}

	if index < 0 || index >= len(data.QAList) {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "索引超出範圍"))
		return
	}

	qa := data.QAList[index]
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("問題: %s\n回答: %s", qa.Question, qa.Answer))
	utils.Send(bot, msg)
}

func delQAInfo(index int) (string, error) {
	if index < 0 || index >= len(data.QAList) {
		return "", fmt.Errorf("index out of range")
	}
	questionStr := data.QAList[index].Question
	data.QAList = append(data.QAList[:index], data.QAList[index+1:]...)
	return questionStr, nil
}

func showDelQA(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var buttons []tgbotapi.InlineKeyboardButton
	for i, qa := range data.QAList {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", i+1, qa.Question), fmt.Sprintf("%s_%d", constant.CALLBACK_QA_DEL, i)))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "請選擇你想刪除的問題:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	utils.Send(bot, msg)
}

func handleQADelButtonClick(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery
	callbackData := cb.Data
	parts := strings.Split(callbackData, "_")
	if len(parts) != 2 {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的指令"))
		return
	}

	command, indexStr := parts[0], parts[1]
	if command != constant.CALLBACK_QA_DEL {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的操作"))
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "無效的索引"))
		return
	}

	if index < 0 || index >= len(data.QAList) {
		utils.Send(bot, tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "索引超出範圍"))
		return
	}

	qa, _ := delQAInfo(index)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("刪除問題: %s", qa))
	utils.Send(bot, msg)
}
