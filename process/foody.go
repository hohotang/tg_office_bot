package process

import (
	"fmt"
	"tgbot/constant"
	"tgbot/data"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleAddingMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	arState := getUserState(message.From.ID)
	switch arState.State {
	case constant.ADD_STATE_NONE:
		AddingFailed(bot, update.Message.From.ID, "狀態錯誤")
		return
	case constant.ADD_STATE_NAME:
		arState.Name = message.Text
		if isRestaurantExist(arState.Name) {
			AddingFailed(bot, update.Message.From.ID, fmt.Sprintf("餐廳%s已經存在", arState.Name))
			return
		}
		arState.State = constant.ADD_STATE_PRICE
		msg := tgbotapi.NewMessage(message.Chat.ID, "請選擇價位:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("低", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_LOW)),
				tgbotapi.NewInlineKeyboardButtonData("中", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_MID)),
				tgbotapi.NewInlineKeyboardButtonData("高", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_HIGH)),
			),
		)
		utils.Send(bot, msg)
	case constant.ADD_STATE_PRICE:
		AddingFailed(bot, update.Message.From.ID, "請使用選項")
	case constant.ADD_STATE_DESCR:
		arState.Description = message.Text
		arState.State = constant.ADD_STATE_CHECK
		// 確認餐廳信息
		confirmationMessage := fmt.Sprintf(
			"請確認以下餐廳信息:\n\n餐廳名稱: %s\n價位: %s\n描述: %s\n\n請輸入 '確認' 或 '取消'",
			arState.Name, constant.PriceStrMap[arState.PriceStr], arState.Description,
		)
		msg := tgbotapi.NewMessage(message.Chat.ID, confirmationMessage)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("確認", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_CONFIRM)),
				tgbotapi.NewInlineKeyboardButtonData("取消", fmt.Sprintf("%s_%s", constant.CALLBACK_ADD, constant.CALLBACK_DENY)),
			),
		)
		utils.Send(bot, msg)
	case constant.ADD_STATE_CHECK:
		AddingFailed(bot, update.Message.From.ID, "請使用選項")
	}
}

func AddingSuccess(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery
	userId := cb.From.ID
	arState := getUserState(userId)
	priceLevel := constant.PriceLevelMap[arState.PriceStr]
	restaurantName := arState.Name
	info := arState.Description
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

	delete(data.UserStates, userID)
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

func IsUserAddingRestaurant(id int64) bool {
	state, exist := data.UserStates[id]
	if !exist {
		return false
	}
	return state.State != constant.ADD_STATE_NONE
}

func getUserState(id int64) *data.AddRestaurantState {
	return data.UserStates[id]
}
