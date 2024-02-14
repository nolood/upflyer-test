package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nolood/upflyer-test.git/internal/config"
	"github.com/nolood/upflyer-test.git/packages/bot"
	"github.com/nolood/upflyer-test.git/packages/parser"
	"github.com/nolood/upflyer-test.git/packages/xlsx"
	"github.com/nolood/upflyer-test.git/utils"
)

func main() {
	config.InitConfig()
	config.InitLogger()
	bot, updates := bot.InitBot()

	// При масштабировании может потребоваться рефакторинг
	for update := range updates {

		if update.Message == nil {
			continue
		}

		var channelName string

		if utils.IsValidChannelURL(update.Message.Text) {
			channelName = utils.ExtractChannelName(update.Message.Text)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, отправьте действительную ссылку на канал в формате https://t.me/channelname")
			bot.Send(msg)
		}

		today := time.Now()
		date := parser.DateOfFirstPost(channelName)
		subs := parser.SubscribersCount(channelName)
		viewsCount, views, err := parser.ViewsAndCountFromLastWeekPosts(channelName)
		if err != nil {
			config.Logger.Error(err.Error())
		}

		medianOfViews := utils.Median(views)
		floatEr := (float64(medianOfViews) / float64(subs)) * 100
		er := math.Round(floatEr*100) / 100
		erString := strconv.FormatFloat(er, 'f', -1, 64)
		age := int(today.Sub(date).Hours() / 24)

		msgText := fmt.Sprintf("DOB %s Age %d \nSubs %d Posts7 %d \nViewsM7 %d \nER %s %s", date.Format("2006.01.02"), age, subs, viewsCount, medianOfViews, erString, "%")

		xlsx.AddToSheet(xlsx.Sheet{Link: update.Message.Text, DOB: date.Format("2006.01.02"), Age: age, Subs: subs, ViewsM7: medianOfViews, ER: erString})

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
