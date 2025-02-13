package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tgbot/conf"
	"tgbot/constant"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProcessCallbackQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery

	if utils.IsAskRestaurantCallback(cb.Data) {
		ProcessAskRestaurantCB(bot, update)
	} else if utils.IsUpdateGameSettingCallback(cb.Data) {
		ProcessUpdateGameSettingCB(bot, update)
	} else if utils.IsQACallback(cb.Data) {
		handleQAButtonClick(bot, update)
	} else if utils.IsQADelCallback(cb.Data) {
		handleQADelButtonClick(bot, update)
	}
}

func activateBat(branchStr string, userName string) {
	config := conf.GetInstance()
	batFilePath := filepath.Join(config.BatPath, config.BatName)
	args := []string{config.StoreFilePath, config.GitServerPath, branchStr, userName}

	if _, err := os.Stat(batFilePath); os.IsNotExist(err) {
		log.Fatalf("Batch file does not exist: %s", batFilePath)
	}

	cmd := exec.Command(batFilePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %v\nOutput:\n%s", err, string(output))
	}

	log.Printf("Command output:\n%s\n", output)
}

func ProcessUpdateGameSettingCB(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery
	var (
		text      string
		groupID   int64
		branchStr string
		actBat    bool
	)
	switch cb.Data {
	case constant.CALL_BACK_TEST:
		text = fmt.Sprintf("開始更新TEST企劃資料 by %s", cb.From.FirstName)
		groupID = conf.GetInstance().TestGroupId
		branchStr = constant.BRANCH_TEST
	default:
		return
	}

	utils.SendMessage(bot, groupID, text)
	if actBat {
		go activateBat(branchStr, cb.From.FirstName)
	}
}

func ProcessAskRestaurantCB(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cb := update.CallbackQuery
	askType, priceLevel := utils.AnalysisRestaurantCallback(cb.Data)
	restaurants := getRestaurants(askType, priceLevel)
	if len(restaurants) == 0 {
		utils.SendMessage(bot, cb.From.ID, "沒有找到任何餐廳")
		return
	}
	var builder strings.Builder
	switch askType {
	case constant.CALLBACK_ALL:
		builder.WriteString("所有餐廳:")
	case constant.CALLBACK_RAND:
		builder.WriteString("推薦餐廳:")
	}
	builder.WriteString(fmt.Sprintf("\n%s", strings.Join(restaurants, "\n")))
	text := builder.String()
	fmt.Println(text)
	utils.SendMessage(bot, cb.From.ID, text)
}
