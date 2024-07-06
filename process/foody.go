package process

import (
	"fmt"
	"tgbot/constant"
	"tgbot/data"
	"tgbot/interactive"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddRestaurantState struct {
	Name        string
	PriceLevel  string
	Description string
}

func askRestaurantName(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	msg := tgbotapi.NewMessage(utils.GetFromID(update), "請輸入餐廳名稱:")
	utils.Send(bot, msg)
	return true
}

func replyRestaurantName(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	state, ok := (*arDataPtr).(*AddRestaurantState)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析餐廳狀態")
		return false
	}
	if isRestaurantExist(update.Message.Text) {
		AddingFailed(bot, utils.GetFromID(update), fmt.Sprintf("餐廳%s已經存在", update.Message.Text))
		return false
	}
	state.Name = update.Message.Text
	return true
}

func askPriceLevel(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	msg := tgbotapi.NewMessage(utils.GetFromID(update), "請選擇價位:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("低", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_LOW)),
			tgbotapi.NewInlineKeyboardButtonData("中", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_MID)),
			tgbotapi.NewInlineKeyboardButtonData("高", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_HIGH)),
		),
	)
	utils.Send(bot, msg)
	return true
}

func replyPriceLevel(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	if update.CallbackQuery == nil {
		AddingFailed(bot, utils.GetFromID(update), "請使用選項")
		return false
	}
	state, ok := (*arDataPtr).(*AddRestaurantState)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析餐廳狀態")
		return false
	}

	callbackData := update.CallbackQuery.Data
	str := utils.AnalysisAddRestaurantCB(callbackData)
	if !processAddPriceState(state, str) {
		AddingFailed(bot, utils.GetFromID(update), "請正確使用選項")
		return false
	}

	return true
}

func askDescription(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	msg := tgbotapi.NewMessage(utils.GetFromID(update), "請輸入餐廳描述:")
	utils.Send(bot, msg)
	return true
}

func replyDescription(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	state := (*arDataPtr).(*AddRestaurantState)
	if update.Message == nil {
		AddingFailed(bot, utils.GetFromID(update), "錯誤")
		return false
	}
	state.Description = update.Message.Text
	*arDataPtr = state
	return true
}

func askConfirm(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	state, ok := (*arDataPtr).(*AddRestaurantState)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析餐廳狀態")
		return false
	}
	// 確認餐廳信息
	confirmationMessage := fmt.Sprintf(
		"請確認以下餐廳信息:\n\n餐廳名稱: %s\n價位: %s\n描述: %s\n\n請輸入 '確認' 或 '取消'",
		state.Name, constant.PriceStrMap[state.PriceLevel], state.Description,
	)
	msg := tgbotapi.NewMessage(utils.GetFromID(update), confirmationMessage)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("確認", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_CONFIRM)),
			tgbotapi.NewInlineKeyboardButtonData("取消", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_DENY)),
		),
	)
	utils.Send(bot, msg)
	return true
}

func replyConfirm(bot *tgbotapi.BotAPI, update *tgbotapi.Update, arDataPtr *interface{}) bool {
	state, ok := (*arDataPtr).(*AddRestaurantState)
	if !ok {
		AddingFailed(bot, utils.GetFromID(update), "內部錯誤，無法解析餐廳狀態")
		return false
	}

	if update.CallbackQuery == nil {
		AddingFailed(bot, utils.GetFromID(update), "請使用選項")
	}
	callbackData := update.CallbackQuery.Data
	str := utils.AnalysisAddRestaurantCB(callbackData)
	if str == constant.CALLBACK_CONFIRM {
		AddingSuccess(bot, update, state)
	} else {
		AddingFailed(bot, utils.GetFromID(update), "新增餐廳已取消")
	}

	return true
}

func registAddRestaurant(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	data := &AddRestaurantState{}
	steps := []*interactive.InteractiveStep{
		{Ask: askRestaurantName, Reply: replyRestaurantName},
		{Ask: askPriceLevel, Reply: replyPriceLevel},
		{Ask: askDescription, Reply: replyDescription},
		{Ask: askConfirm, Reply: replyConfirm},
	}
	interactive.RegisInteractiveMode(update.Message.From.ID, data, steps)
	steps[0].Ask(bot, update, nil)
}

func AddingSuccess(bot *tgbotapi.BotAPI, update *tgbotapi.Update, state *AddRestaurantState) {
	cb := update.CallbackQuery
	userId := cb.From.ID
	priceLevel := constant.PriceLevelMap[state.PriceLevel]
	restaurantName := state.Name
	info := state.Description
	userName := fmt.Sprintf("%s %s", cb.From.FirstName, cb.From.LastName)
	data.RestaurantMap[priceLevel][restaurantName] = data.RestaurantInfo{
		Recommender: userName,
		RecID:       userId,
		Info:        info,
	}
	sendMessageAndClearState(bot, userId, fmt.Sprintf("感謝 %s, 成功加入餐廳: %s", userName, restaurantName))

}

func AddingFailed(bot *tgbotapi.BotAPI, userID int64, errorMsg string) {
	sendMessageAndClearState(bot, userID, fmt.Sprintf("*新增餐廳錯誤: %s*", errorMsg))
}

func sendMessageAndClearState(bot *tgbotapi.BotAPI, userID int64, message string) {
	utils.SendMessage(bot, userID, message)
	interactive.DelInteractive(userID)
}

func sendPriceLevelOptions(bot *tgbotapi.BotAPI, chatID int64, commandType string) {
	msg := tgbotapi.NewMessage(chatID, "請選擇價位")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("低", fmt.Sprintf("%s_%s_%s", constant.CALLBACK_RESTAURANT, commandType, constant.CALLBACK_LOW)),
			tgbotapi.NewInlineKeyboardButtonData("中", fmt.Sprintf("%s_%s_%s", constant.CALLBACK_RESTAURANT, commandType, constant.CALLBACK_MID)),
			tgbotapi.NewInlineKeyboardButtonData("高", fmt.Sprintf("%s_%s_%s", constant.CALLBACK_RESTAURANT, commandType, constant.CALLBACK_HIGH)),
		),
	)
	utils.Send(bot, msg)
}

// Utility functions

func deleteRestaurant(restaurantName string) bool {
	for _, rm := range data.RestaurantMap {
		if _, exist := rm[restaurantName]; exist {
			delete(rm, restaurantName)
			return true
		}
	}
	return false
}

func getRestaurants(askType, priceLevel string) []string {
	level, ok := constant.PriceLevelMap[priceLevel]
	if !ok {
		return []string{}
	}

	switch askType {
	case constant.CALLBACK_ALL:
		return getAllRestaurants(level)
	case constant.CALLBACK_RAND:
		return getRandomRestaurants(level)
	}
	return []string{}
}

func getRandomRestaurants(priceLevel int) []string {
	return formatRestaurants(priceLevel, true, constant.RECOMMEND_NUM)
}

func getAllRestaurants(priceLevel int) []string {
	return formatRestaurants(priceLevel, false, 0)
}

func formatRestaurants(priceLevel int, shuffle bool, limit int) []string {
	rm := data.RestaurantMap[priceLevel]
	rNameSlice := make([]string, 0, len(rm))
	for name := range rm {
		rNameSlice = append(rNameSlice, name)
	}

	if shuffle {
		rNameSlice = utils.ShuffleSlice(rNameSlice)
	}

	restaurants := make([]string, 0, len(rNameSlice))
	for i, name := range rNameSlice {
		if limit > 0 && i >= limit {
			break
		}
		info := rm[name]
		restaurants = append(restaurants, fmt.Sprintf("餐廳: %s\n推薦人: %s\n心得 :%s\n---------------------------------", name, info.Recommender, info.Info))
	}
	return restaurants
}

func isRestaurantExist(restaurantName string) bool {
	for _, rm := range data.RestaurantMap {
		_, exist := rm[restaurantName]
		if exist {
			return true
		}
	}
	return false
}

func processAddPriceState(arState *AddRestaurantState, str string) bool {
	if str == constant.CALLBACK_LOW || str == constant.CALLBACK_MID || str == constant.CALLBACK_HIGH {
		arState.PriceLevel = str
	} else {
		return false
	}
	return true
}
