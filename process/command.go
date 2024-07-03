package process

import (
	"fmt"
	"strings"
	"tgbot/conf"
	"tgbot/constant"
	"tgbot/data"
	"tgbot/reminder"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandHelp               = "help"
	commandStart              = "start"
	commandAskExcelPermission = "ask_excel_permission"
	commandUpdateGameSetting  = "update_gamesetting"
	commandAddExcel           = "add_excel"
	commandDelExcel           = "del_excel"
	commandListExcel          = "list_excel"
	commandAddBan             = "add_ban"
	commandDelBan             = "del_ban"
	commandAddFoody           = "add_foody"
	commandDelFoody           = "del_foody"
	commandListFoody          = "list_foody"
	commandAddRestaurant      = "add_restaurant"
	commandDelRestaurant      = "del_restaurant"
	commandAskRandRestaurant  = "ask_rand_restaurant"
	commandAskAllRestaurant   = "ask_all_restaurant"
	commandAskFoodyPermission = "ask_foody_permission"
	commandReminderSwitch     = "reminder_switch"
	commandGetGif             = "get_a_gif"
	commandGetMusic           = "get_music"
	commandReboot             = "reboot"
	commandTest               = "test"

	msgNoSuchCmd               = "No such command!!!"
	msgUnauthorized            = "You are not authorized to perform this action"
	msgInvalidIDFormat         = "內容錯誤:"
	msgInvalidRestaurantFormat = " 請按照正確格式填寫"
	msgReminderOnMessage       = "提醒功能開啟"
	msgReminderOffMessage      = "提醒功能關閉"
	msgInvalidArgsMessage      = "請輸入 on 或 off"
)

// CommandHandler is the interface that all command handlers should implement
type CommandHandler interface {
	HandleCommand(*tgbotapi.BotAPI, *tgbotapi.Update)
}

// CommandRouter routes commands to their handlers
type CommandRouter struct {
	handlers map[string]CommandHandler
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]CommandHandler),
	}
}

func (r *CommandRouter) Register(command string, handler CommandHandler) {
	r.handlers[command] = handler
}

func (r *CommandRouter) Route(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if handler, found := r.handlers[update.Message.Command()]; found {
		handler.HandleCommand(bot, update)
	} else {
		utils.SendMarkDownMessage(bot, update.Message.Chat.ID, "No such command!!!")
	}
}

// private command

type StartCommand struct{}

func (h *StartCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	utils.SendMessage(bot, update.Message.Chat.ID, "安安 這是辦公室機器人")
}

type HelpCommand struct{}

func (h *HelpCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	var builder strings.Builder

	builder.WriteString("可使用指令:")

	if utils.IsAdmin(update.Message.From.ID) {
		appendAdminCommands(&builder)
	}

	if utils.IsUpdateGameSettingLegal(update.Message.From.ID) {
		appendUpdateGameSettingCommands(&builder)
	} else {
		builder.WriteString(fmt.Sprintf("\n/%s : 詢問更新excel權限", commandAskExcelPermission))
	}

	if utils.IsFoody(update.Message.From.ID) {
		appendFoodyCommands(&builder)
	} else {
		builder.WriteString(fmt.Sprintf("\n/%s : 詢問饕客權限", commandAskFoodyPermission))
	}

	appendGeneralCommands(&builder)

	utils.SendMessage(bot, chatID, builder.String())
}

type AskExcelPermissionCommand struct{}

func (h *AskExcelPermissionCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := fmt.Sprintf("%s 詢問更新excel權限 id : %d", update.Message.From.FirstName+update.Message.From.LastName, update.Message.From.ID)
	for adminID := range data.AdminMap {
		utils.SendMessage(bot, adminID, text)
	}
}

type UpdateGameSettingCommand struct{}

func (h *UpdateGameSettingCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsUpdateGameSettingLegal(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	text := fmt.Sprintf("Hello %s\n更新環境:", update.Message.From.FirstName)
	SendGameSettingOption(bot, update.Message.Chat.ID, text)
}

type AddExcelCommand struct{}

func (h *AddExcelCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	chatID := update.Message.Chat.ID
	id, userName, err := utils.ParseIDAndUserNameFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, chatID, msgInvalidIDFormat+err.Error())
		return
	}
	data.ExcelPermissionMap[id] = userName
	utils.SendMessage(bot, chatID, fmt.Sprintf("Added permission for ID %d", id))
	utils.SendMessage(bot, id, "已獲得新增excel及更新權限")
}

type DelExcelCommand struct{}

func (h *DelExcelCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	chatID := update.Message.Chat.ID
	id, err := utils.ParseIDFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, chatID, msgInvalidIDFormat)
		return
	}
	delete(data.ExcelPermissionMap, id)
	utils.SendMessage(bot, chatID, fmt.Sprintf("Deleted permission for ID %d", id))
}

type AddBanCommand struct{}

func (h *AddBanCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	chatID := update.Message.Chat.ID
	id, userName, err := utils.ParseIDAndUserNameFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, chatID, msgInvalidIDFormat)
		return
	}
	data.BannedMap[id] = userName
	utils.SendMessage(bot, chatID, fmt.Sprintf("Banned user with ID %d", id))
}

type DelBanCommand struct{}

func (h *DelBanCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	chatID := update.Message.Chat.ID
	id, err := utils.ParseIDFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, chatID, msgInvalidIDFormat)
		return
	}
	delete(data.BannedMap, id)
	utils.SendMessage(bot, chatID, fmt.Sprintf("Unbanned user with ID %d", id))
}

type AddFoodyCommand struct{}

func (h *AddFoodyCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	chatID := update.Message.Chat.ID
	id, userName, err := utils.ParseIDAndUserNameFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, chatID, msgInvalidIDFormat)
		return
	}
	data.FoodyMap[id] = userName
	utils.SendMessage(bot, chatID, fmt.Sprintf("Added foody with ID %d", id))
	utils.SendMessage(bot, id, "已獲得饕客權限(可以新增餐廳)")
}

type DelFoodyCommand struct{}

func (h *DelFoodyCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	id, err := utils.ParseIDFromCommandArguments(update)
	if err != nil {
		utils.SendMessage(bot, update.Message.Chat.ID, msgInvalidIDFormat)
		return
	}
	delete(data.FoodyMap, id)
	utils.SendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("Deleted foody with ID %d", id))
}

type ListFoodyCommand struct{}

func (h *ListFoodyCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	var msg string
	for id, name := range data.FoodyMap {
		msg += fmt.Sprintf("%d: %s\n", id, name)
	}
	utils.SendMessage(bot, update.Message.Chat.ID, msg)
}

type AddRestaurantCommand struct{}

func (h *AddRestaurantCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsFoody(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	data.UserStates[update.Message.From.ID] = &data.AddRestaurantState{
		State: constant.ADD_STATE_NAME,
	}
	utils.SendMessage(bot, update.Message.From.ID, "請輸入餐廳名稱:")
}

type DelRestaurantCommand struct{}

func (h *DelRestaurantCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	restaurantName := utils.ParseNameFromCommandArguments(update)
	if !utils.CanDeleteRestaurant(update.Message.From.ID, restaurantName) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}
	if deleteRestaurant(restaurantName) {
		utils.SendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("刪除餐廳: %s", restaurantName))
	} else {
		utils.SendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("沒有找到餐廳: %s", restaurantName))
	}
}

type AskRandomRestaurantCommand struct{}

func (h *AskRandomRestaurantCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	sendPriceLevelOptions(bot, update.Message.Chat.ID, constant.CALLBACK_RAND)
}

type AskAllRestaurantCommand struct{}

func (h *AskAllRestaurantCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	sendPriceLevelOptions(bot, update.Message.Chat.ID, constant.CALLBACK_ALL)
}

type AskFoodyPermissionCommand struct{}

func (h *AskFoodyPermissionCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := fmt.Sprintf("%s %s 詢問更新餐廳權限 id : %d", update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.ID)
	for adminID := range data.AdminMap {
		utils.SendMessage(bot, adminID, text)
	}
}

type ReminderSwitchCommand struct{}

func (h *ReminderSwitchCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !utils.IsAdmin(update.Message.From.ID) {
		utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
		return
	}

	nowState := reminder.ReminderSwitch()
	utils.SendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("當前提醒狀態:%v", nowState))
}

type GetGifCommand struct{}

func (h *GetGifCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	gifPath, err := utils.GetRandomGif(conf.GetInstance().MemePath)
	if err != nil {
		utils.SendMessage(bot, chatID, "Error getting gif: "+err.Error())
		return
	}
	utils.SendGif(bot, chatID, gifPath)
}

type GetMusicCommand struct{}

func (h *GetMusicCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	gifPath, err := utils.GetRandomMusic(conf.GetInstance().MusicPath)
	if err != nil {
		utils.SendMessage(bot, chatID, "Error getting music: "+err.Error())
		return
	}
	utils.SendMusic(bot, chatID, gifPath)
}

type RebootCommand struct{}

func (h *RebootCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// if !utils.IsAdmin(update.Message.From.ID) {
	// 	utils.SendMessage(bot, update.Message.Chat.ID, msgUnauthorized)
	// 	return
	// }
	// utils.SendMessage(bot, update.Message.Chat.ID, "Rebooting...")
	// os.Exit(0)
}

type TestCommand struct{}

func (h *TestCommand) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	utils.SaveDataToFile(conf.GetInstance().SaveFileName)
}

func InitializePrivateCommandRouter() *CommandRouter {
	router := NewCommandRouter()

	router.Register(commandStart, &StartCommand{})
	router.Register(commandHelp, &HelpCommand{})
	router.Register(commandAskExcelPermission, &AskExcelPermissionCommand{})
	router.Register(commandUpdateGameSetting, &UpdateGameSettingCommand{})
	router.Register(commandAddExcel, &AddExcelCommand{})
	router.Register(commandDelExcel, &DelExcelCommand{})
	router.Register(commandAddBan, &AddBanCommand{})
	router.Register(commandDelBan, &DelBanCommand{})
	router.Register(commandAddFoody, &AddFoodyCommand{})
	router.Register(commandDelFoody, &DelFoodyCommand{})
	router.Register(commandListFoody, &ListFoodyCommand{})
	router.Register(commandAddRestaurant, &AddRestaurantCommand{})
	router.Register(commandDelRestaurant, &DelRestaurantCommand{})
	router.Register(commandAskRandRestaurant, &AskRandomRestaurantCommand{})
	router.Register(commandAskAllRestaurant, &AskAllRestaurantCommand{})
	router.Register(commandAskFoodyPermission, &AskFoodyPermissionCommand{})
	router.Register(commandReminderSwitch, &ReminderSwitchCommand{})
	router.Register(commandGetGif, &GetGifCommand{})
	router.Register(commandGetMusic, &GetMusicCommand{})
	router.Register(commandReboot, &RebootCommand{})
	router.Register(commandTest, &TestCommand{})

	return router
}

// public command

type PublicStartCmd struct{}

func (h *PublicStartCmd) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	sendChatResponse(bot, update.Message.Chat.ID)
}

type PublicHelpCmd struct{}

func (h *PublicHelpCmd) HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	sendChatResponse(bot, update.Message.Chat.ID)
}

func InitializePublicCommandRouter() *CommandRouter {
	router := NewCommandRouter()

	router.Register(commandStart, &PublicStartCmd{})
	router.Register(commandHelp, &PublicHelpCmd{})

	return router
}

// sendChatResponse sends a default response for chat commands.
func sendChatResponse(bot *tgbotapi.BotAPI, chatID int64) {
	utils.SendMessage(bot, chatID, fmt.Sprintf("請跟我對話 @%s", bot.Self.UserName))
}

// SendGameSettingOption sends game setting options to the user.
func SendGameSettingOption(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("TEST", constant.CALL_BACK_TEST),
		),
	)
	utils.Send(bot, msg)
}

// appendAdminCommands appends admin commands to the string builder.
func appendAdminCommands(builder *strings.Builder) {
	builder.WriteString("\n-------管理員權限-------")
	builder.WriteString(fmt.Sprintf("\n/%s : 新增excel權限 ", commandAddExcel))
	builder.WriteString(fmt.Sprintf("\n/%s : 刪除excel權限", commandDelExcel))
	builder.WriteString(fmt.Sprintf("\n/%s : 新增饕客權限", commandAddFoody))
	builder.WriteString(fmt.Sprintf("\n/%s : 刪除饕客權限", commandDelFoody))
	builder.WriteString(fmt.Sprintf("\n/%s : 顯示所有饕客", commandListFoody))
	builder.WriteString(fmt.Sprintf("\n/%s : 新增禁止使用者", commandAddBan))
	builder.WriteString(fmt.Sprintf("\n/%s : 刪除禁止使用者", commandDelBan))
	builder.WriteString(fmt.Sprintf("\n/%s : 切換提醒開關", commandReminderSwitch))
	builder.WriteString("\n-------管理員權限-------")
}

// appendUpdateGameSettingCommands appends update gamesetting permission commands to the string builder.
func appendUpdateGameSettingCommands(builder *strings.Builder) {
	builder.WriteString("\n-------Excel權限-------")
	builder.WriteString(fmt.Sprintf("\n/%s : 更新GameSetting，請先儲存GameSetting.xlsx\n儲存方式:拖曳GameSetting.xlsx到群組即可", commandUpdateGameSetting))
	builder.WriteString("\n-------Excel權限-------")
}

// appendFoodyCommands appends foody commands to the string builder.
func appendFoodyCommands(builder *strings.Builder) {
	builder.WriteString("\n-------饕客權限-------")
	builder.WriteString(fmt.Sprintf("\n/%s : 新增餐廳 (建議推薦辦公室附近的餐廳)", commandAddRestaurant))
	builder.WriteString(fmt.Sprintf("\n/%s : 刪除餐廳 只能刪除自己新增的餐廳喔", commandDelRestaurant)) // TODO 用選項讓使用者刪除
	builder.WriteString(fmt.Sprintf("\n例如: /%s 八方雲集", commandDelRestaurant))
	builder.WriteString("\n-------饕客權限-------")
}

// appendGeneralCommands appends general commands to the string builder.
func appendGeneralCommands(builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("\n/%s : 快樂一下", commandGetGif))
	builder.WriteString(fmt.Sprintf("\n/%s : 隨機一首歌", commandGetMusic))
	builder.WriteString(fmt.Sprintf("\n/%s : 推薦三家符合價位的餐廳", commandAskRandRestaurant))
	builder.WriteString(fmt.Sprintf("\n/%s : 找尋所有符合價位的餐廳", commandAskAllRestaurant))
}

// SendBannedMessage sends a message to a banned user.
func SendBannedMessage(bot *tgbotapi.BotAPI, chatID int64) {
	utils.SendMessage(bot, chatID, "You are banned from using this bot")
}
