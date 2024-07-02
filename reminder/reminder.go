package reminder

import (
	"log"
	"tgbot/conf"
	"tgbot/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

const (
	reminder_good_morning = "good_morning"
	reminder_wanting_eat  = "wanting_eat"
	reminder_weekend      = "weekend"
	reminder_get_off      = "get_off"
	reminder_test         = "test"
)

type ReminderConfig struct {
	Name    string
	Spec    string
	ChatID  int64
	Message string
}

var reminderSwitch bool

func Start(bot *tgbotapi.BotAPI) {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Panic(err)
	}
	c := cron.New(cron.WithLocation(loc))
	reminderSwitch = true

	reminders := []ReminderConfig{
		{reminder_test, "30 14 * * 1-5", conf.TEST_GROUP_ID, "*測試*"},
	}

	for _, reminder := range reminders {
		addReminder(c, bot, reminder)
	}

	// 启动定时任务
	c.Start()

	// 防止主协程退出
	select {}
}

func addReminder(c *cron.Cron, bot *tgbotapi.BotAPI, reminder ReminderConfig) {
	_, err := c.AddFunc(reminder.Spec, func() {
		if reminderSwitch {
			utils.SendMarkDownMessage(bot, reminder.ChatID, reminder.Message)
			switch reminder.Name {
			case reminder_good_morning:
				SendGoodMorningGif(bot, reminder.ChatID)
			case reminder_weekend:
				SendWeekendGif(bot, reminder.ChatID)
			case reminder_test:
				SendWeekendGif(bot, conf.TEST_GROUP_ID)
			}
		}
	})
	if err != nil {
		log.Panic(err)
	}
}

func ReminderSwitch() bool {
	reminderSwitch = !reminderSwitch
	return reminderSwitch
}

func SendGoodMorningGif(bot *tgbotapi.BotAPI, chatID int64) {
	SendGif(bot, chatID, conf.GOOD_MORING_PATH)
}

func SendWeekendGif(bot *tgbotapi.BotAPI, chatID int64) {
	SendGif(bot, chatID, conf.WEEKEND_PATH)
}

func SendGif(bot *tgbotapi.BotAPI, chatID int64, gifSource string) {
	gifPath, err := utils.GetRandomGif(gifSource)
	if err != nil {
		utils.SendMessage(bot, chatID, "Error getting gif: "+err.Error())
		return
	}
	utils.SendGif(bot, chatID, gifPath)
}
