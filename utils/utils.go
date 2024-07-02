package utils

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tgbot/conf"
	"tgbot/constant"
	"tgbot/data"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Send(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	DelLastMessage(bot, msg.ChatID)
	message, _ := bot.Send(msg)
	StoreLastMessage(bot, msg.ChatID, message.MessageID)
}

// SendMessage sends a command message to a specific chat. with no parse problem like "_"
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	Send(bot, msg)
}

// SendMarkDownMessage sends a text message to a specific chat.
func SendMarkDownMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, _ = bot.Send(msg)
}

func StoreLastMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	data.LastMessageID[chatID] = messageID
}

func IsChatPrivate(update *tgbotapi.Update) bool {
	return update.Message.Chat.Type == "private"
}

func DelLastMessage(bot *tgbotapi.BotAPI, chatID int64) {
	messageID, exist := data.LastMessageID[chatID]
	if !exist {
		return
	}
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = bot.Send(msg)
}

func GetRandomGif(directory string) (string, error) {
	files, err := filepath.Glob(filepath.Join(directory, "*.gif"))
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", os.ErrNotExist
	}

	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(files))
	return files[randomIndex], nil
}

func SendGif(bot *tgbotapi.BotAPI, chatID int64, gifPath string) {
	fileName := filepath.Base(gifPath)
	if IsFileExist(fileName) {
		err := SendGifFileID(bot, chatID, data.Cache[fileName])
		if err == nil {
			return
		}
	}
	gif := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(gifPath))
	message, err := bot.Send(gif)
	if err != nil {
		panic(err)
	}
	if message.Document != nil {
		if _, exist := data.Cache[fileName]; !exist {
			data.Cache[fileName] = message.Document.FileID
		}
	}
}

func GetRandomMusic(directory string) (string, error) {
	extensions := []string{"*.wav", "*.mp3", "*.flac"}
	var files []string

	for _, ext := range extensions {
		matchedFiles, err := filepath.Glob(filepath.Join(directory, ext))
		if err != nil {
			return "", err
		}
		files = append(files, matchedFiles...)
	}

	if len(files) == 0 {
		return "", os.ErrNotExist
	}

	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(files))
	return files[randomIndex], nil
}

func SendMusic(bot *tgbotapi.BotAPI, chatID int64, musicPath string) {
	fileName := filepath.Base(musicPath)
	if IsFileExist(fileName) {
		err := SendMusicFileID(bot, chatID, data.Cache[fileName])
		if err == nil {
			return
		}
	}
	gif := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(musicPath))
	message, err := bot.Send(gif)
	if err != nil {
		panic(err)
	}
	if message.Document != nil {
		if _, exist := data.Cache[fileName]; !exist {
			data.Cache[fileName] = message.Document.FileID
		}
	}
	if message.Audio != nil {
		if _, exist := data.Cache[fileName]; !exist {
			data.Cache[fileName] = message.Audio.FileID
		}
	}
}

func IsFileExist(fileName string) bool {
	return data.Cache[fileName] != ""
}

func SendGifFileID(bot *tgbotapi.BotAPI, chatID int64, fileID string) error {
	gif := tgbotapi.NewDocument(chatID, tgbotapi.FileID(fileID))
	_, err := bot.Send(gif)
	return err
}

func SendMusicFileID(bot *tgbotapi.BotAPI, chatID int64, fileID string) error {
	music := tgbotapi.NewAudio(chatID, tgbotapi.FileID(fileID))
	_, err := bot.Send(music)
	return err
}

// IsUpdateGameSettingLegal checks if a user has permission to update game settings.
func IsUpdateGameSettingLegal(id int64) bool {
	_, exist := data.ExcelPermissionMap[id]
	return exist
}

// IsAdmin checks if a user is an admin.
func IsAdmin(id int64) bool {
	_, exist := data.AdminMap[id]
	return exist
}

// IsBanned checks if a user is banned.
func IsBanned(id int64) bool {
	_, exist := data.BannedMap[id]
	return exist
}

// IsAdmin checks if a user is a foody.
func IsFoody(id int64) bool {
	_, exist := data.FoodyMap[id]
	return exist
}

func CanDeleteRestaurant(id int64, restaurantName string) bool {
	if IsAdmin(id) || IsRecommender(id, restaurantName) {
		return true
	}
	return false
}

func IsRecommender(id int64, restaurantName string) bool {
	for _, rm := range data.RestaurantMap {
		restaurant, exist := rm[restaurantName]
		if !exist {
			continue
		}
		if restaurant.RecID == id {
			return true
		}
	}
	return false
}

func ShuffleSlice[T any](obj []T) []T {
	rand.Shuffle(len(obj), func(i, j int) {
		obj[i], obj[j] = obj[j], obj[i]
	})
	return obj
}

func SaveDataToFile(filename string) error {
	data := data.SaveData{
		ExcelPermissionMap: data.ExcelPermissionMap,
		AdminMap:           data.AdminMap,
		BannedMap:          data.BannedMap,
		Cache:              data.Cache,
		FoodyMap:           data.FoodyMap,
		RestaurantMap:      data.RestaurantMap[:],
	}

	file, err := os.Create(conf.SAVE_FILE_PATH + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // for pretty printing
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func LoadDataFromFile(filename string) error {
	file, err := os.Open(conf.SAVE_FILE_PATH + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	saveData := data.SaveData{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&saveData); err != nil {
		return err
	}

	// Merge the loaded data with the existing maps
	for k, v := range saveData.ExcelPermissionMap {
		data.ExcelPermissionMap[k] = v
	}
	for k, v := range saveData.AdminMap {
		data.AdminMap[k] = v
	}
	for k, v := range saveData.BannedMap {
		data.BannedMap[k] = v
	}
	for k, v := range saveData.Cache {
		data.Cache[k] = v
	}
	for k, v := range saveData.FoodyMap {
		data.FoodyMap[k] = v
	}
	for pl, rm := range saveData.RestaurantMap {
		for k, v := range rm {
			data.RestaurantMap[pl][k] = v
		}
	}

	return nil
}

// ParseIDFromCommandArguments parses an ID from the command arguments.
func ParseIDFromCommandArguments(update *tgbotapi.Update) (int64, error) {
	args := update.Message.CommandArguments()

	id, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// parseIDFromCommandArguments parses an ID and an userName from the command arguments.
func ParseIDAndUserNameFromCommandArguments(update *tgbotapi.Update) (int64, string, error) {
	args := update.Message.CommandArguments()
	parts := strings.Fields(args)

	if len(parts) < 2 {
		return 0, "", errors.New("參數不夠")
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", err
	}

	userName := strings.Join(parts[1:], " ")

	return id, userName, nil
}

func ParseNameFromCommandArguments(update *tgbotapi.Update) string {
	return update.Message.CommandArguments()
}

func IsUpdateGameSettingCallback(cs string) bool {
	parts := strings.Split(cs, "_")
	return parts[0] == constant.CALLBACK_UPDATE
}

func IsAskRestaurantCallback(cs string) bool {
	parts := strings.Split(cs, "_")
	return parts[0] == constant.CALLBACK_RESTAURANT
}

func IsAddRestaurantCallback(cs string) bool {
	parts := strings.Split(cs, "_")
	return parts[0] == constant.CALLBACK_ADD
}

func AnalysisRestaurantCallback(cs string) (string, string) {
	parts := strings.Split(cs, "_")
	askType := parts[1]
	priceType := parts[2]
	return askType, priceType
}

func AnalysisAddRestaurantCB(cs string) string {
	parts := strings.Split(cs, "_")
	return parts[1]
}
