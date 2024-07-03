package app

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"tgbot/conf"
	"tgbot/process"
	"tgbot/reminder"
	"tgbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DefaultApp struct {
	version              string
	processID            string
	privateCommandRouter *process.CommandRouter
	publicCommandRouter  *process.CommandRouter
	mu                   sync.Mutex
}

func NewApp(version string) *DefaultApp {
	return &DefaultApp{
		processID:            "development",
		version:              version,
		privateCommandRouter: process.InitializePrivateCommandRouter(),
		publicCommandRouter:  process.InitializePublicCommandRouter(),
	}
}

func (app *DefaultApp) Run(debug bool) {
	defer recoverPanic()

	config := conf.GetInstance()

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debug

	if err := utils.LoadDataFromFile(config.SaveFileName); err != nil {
		log.Printf("Failed to save data to file: %v", err)
	}

	// Set up channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Set up a channel to notify the main goroutine to exit
	doneChan := make(chan bool, 1)
	// Signal handler goroutine
	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %s", sig)
		// Save data to file
		err = os.MkdirAll(config.SaveFilePath, os.ModePerm)
		if err := utils.SaveDataToFile(config.SaveFileName); err != nil {
			log.Printf("Failed to save data to file: %v", err)
		}
		doneChan <- true
	}()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	go reminder.Start(bot)

	for {
		select {
		case update := <-updates:
			app.handleUpdate(bot, &update)
		case <-doneChan:
			log.Println("Shutting down gracefully...")
			return
		}
	}
}

func (app *DefaultApp) handleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	defer func() {
		app.mu.Unlock()
		recoverPanic()
	}()

	app.mu.Lock()
	if update.Message != nil {
		if utils.IsBanned(update.Message.From.ID) {
			process.SendBannedMessage(bot, update.Message.From.ID)
			return
		}
		if process.IsUserAddingRestaurant(update.Message.From.ID) {
			process.HandleAddingMessage(bot, update)
		} else if update.Message.IsCommand() && utils.IsChatPrivate(update) {
			app.privateCommandRouter.Route(bot, update)
		} else if update.Message.IsCommand() && !utils.IsChatPrivate(update) {
			app.publicCommandRouter.Route(bot, update)
		} else if update.Message.Document != nil {
			process.ProcessDocument(bot, update)
		}
	}

	if update.CallbackQuery != nil {
		process.ProcessCallbackQuery(bot, update)
	}
}

func recoverPanic() {
	if r := recover(); r != nil {
		var errMsg string
		switch v := r.(type) {
		case string:
			errMsg = v
		case error:
			errMsg = v.Error()
		default:
			errMsg = "Unknown panic"
		}
		buf := make([]byte, 1024)
		length := runtime.Stack(buf, false)
		stackTrace := string(buf[:length])
		log.Printf("Recovered from panic: %s\nStack trace:\n%s", errMsg, stackTrace)
	}
}
