package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nolood/upflyer-test.git/internal/config"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

func main() {
	config.InitConfig()
	config.InitLogger()

	token := viper.GetString("TELEGRAM_BOT_TOKEN")
	if token == "" {
		config.Logger.Panic("TELEGRAM_BOT_TOKEN not found")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//channelRegex := regexp.MustCompile(`t(?:elegram)?\.me\/([^\/]+)\/?`)

	//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	//matches := channelRegex.FindStringSubmatch(update.Message.Text)

	//if len(matches) > 1 {
	//channelUsername := matches[1]

	//}
	//bot.Send(msg)

	for update := range updates {
		if update.Message != nil { // If we got a message

			channelID, err := getChannelID(update.Message.Text)
			if err != nil {
				config.Logger.Error(err.Error())
				return
			}

			channel, err := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: channelID}})
			if err != nil {
				config.Logger.Error(err.Error())
				continue
			}

			data, err := getFirstPost(strings.Split(update.Message.Text, "/")[3])
			if err != nil {
				config.Logger.Error(err.Error())
				continue
			}

			log.Println(data, "data------------------")
			//firstMessage(bot)

			//chatHistory(token, strings.Split(update.Message.Text, "/")[3]

			//log.Println(channel.Description)

			msgText := fmt.Sprintf("Earliest data %s: \n Title: %s", "ok", channel.Title)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

			bot.Send(msg)
		}
	}
}

func getFirstPost(channelUsername string) (string, error) {

	log.Println("getFirstPost", channelUsername)
	url := fmt.Sprintf("https://t.me/s/%s", channelUsername)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}

	postDiv := doc.Find(".tgme_page_post")
	if postDiv.Length() > 0 {
		postText := postDiv.Find(".tgme_page_post_body").Text()
		return strings.TrimSpace(postText), nil
	}

	return "", fmt.Errorf("no post found")
}

//type MyParams struct {
//	Peer       int64
//	Limit      int
//	OffsetID   int
//	OffsetDate int
//	AddOffset  int
//	MaxID      int
//	MinID      int
//}

//func firstMessage(bot *tgbotapi.BotAPI, channelID int64) {
//
//	params := MyParams{
//		Peer:       channelID,
//		Limit:      1,
//		OffsetID:   0,
//		OffsetDate: 0,
//		AddOffset:  0,
//		MaxID:      0,
//		MinID:      0,
//	}
//
//	data, err := bot.MakeRequest("messages.getHistory")
//	if err != nil {
//		log.Println(err)
//	}
//
//	log.Println(data, "----------------data")
//}

//func firstMessage(channelID int64) {
//	client := telegram.NewClient(telegram.TestAppID, telegram.TestAppHash, telegram.Options{
//		//Logger: config.Logger,
//	})
//
//	ctx := context.Background()
//	if err := client.Run(ctx, func(ctx context.Context) error {
//		peer := tg.InputPeerChannel{
//			ChannelID:  channelID,
//			AccessHash: 0, // Если канал публичный, то AccessHash равен 0
//		}
//
//		resp, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
//			Peer:       &peer,
//			Limit:      1,
//			OffsetID:   0,
//			OffsetDate: 0,
//			AddOffset:  0,
//			MaxID:      0,
//			MinID:      0,
//		})
//		if err != nil {
//			return err
//		}
//
//		log.Println(resp, "resp------------------------------------+")
//
//		return nil
//	}); err != nil {
//		//config.Logger.Error(err.Error())
//	}
//}

//type message struct {
//	MessageID int    `json:"message_id"`
//	Text      string `json:"text"`
//}
//
//type messageResponse struct {
//	Result []message `json:"result"`
//}

//func chatHistory(token string, channelUsername string) {
//
//	log.Println("chatHistory", channelUsername)
//
//	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChatHistory?chat_id=@%s&limit=1&chat_type=channel", token, channelUsername)
//
//	// Отправляем GET запрос к API Telegram
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	// Читаем и декодируем JSON ответ
//	var messages messageResponse
//	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println(messages, "messages")
//
//	// Если есть сообщение в канале, выводим его
//	if len(messages.Result) > 0 {
//		firstMessage := messages.Result[0]
//		log.Printf("First message from the channel: %v", firstMessage.Text)
//	} else {
//		log.Println("No messages found in the channel.")
//	}
//}

type GetChatResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		Id int64 `json:"id"`
	} `json:"result"`
}

func getChannelID(channelURL string) (int64, error) {
	parts := strings.Split(channelURL, "/")
	botToken := viper.GetString("TELEGRAM_BOT_TOKEN")
	channelUsername := parts[len(parts)-1]

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=@%s",
		botToken, channelUsername)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var chatResponse GetChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return 0, err
	}

	if chatResponse.Ok {
		return chatResponse.Result.Id, nil
	}

	return 0, fmt.Errorf("channel ID not found")
}
