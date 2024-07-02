package process

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"tgbot/conf"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsGameSetting(fileName string) bool {
	match, _ := regexp.MatchString(`^GameSetting\.xlsx$`, fileName)
	return match
}

func ProcessDocument(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message.Chat.Type != "private" {
		return
	}
	if utils.IsUpdateGameSettingLegal(update.Message.From.ID) {
		if IsGameSetting(update.Message.Document.FileName) {
			err := ProcessExcel(bot, update)
			if err != nil {
				log.Printf("Failed to process Excel file: %v", err)
			}
		} else {
			utils.SendMessage(bot, update.Message.Chat.ID, "不是GameSetting.xlsx")
		}
	} else {
		utils.SendMessage(bot, update.Message.Chat.ID, "有權限才能更新喔")
	}
}

func ProcessExcel(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	fileID := update.Message.Document.FileID
	fileURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Printf("Failed to get file URL: %v", err)
		return err
	}

	// 下载文件
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Printf("Failed to download file: %v", err)
		return err
	}
	defer resp.Body.Close()

	saveDir := conf.STORE_FILE_PATH
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directory: %v", err)
		return err
	}

	filePath := filepath.Join(saveDir, update.Message.Document.FileName)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		return err
	}

	utils.SendMessage(bot, update.Message.Chat.ID, fmt.Sprintf("文件%s已經成功保存", update.Message.Document.FileName))
	return nil
}
