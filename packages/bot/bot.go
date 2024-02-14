package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nolood/upflyer-test.git/internal/config"
	"github.com/spf13/viper"
)

func InitBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {

	token := viper.GetString("TELEGRAM_BOT_TOKEN")
	if token == "" {
		config.Logger.Panic("TELEGRAM_BOT_TOKEN not found")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return bot, updates

}
