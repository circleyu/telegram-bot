package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramBot app
type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramBot function
func NewTelegramBot(apiToken string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		bot: bot,
	}, nil
}

func init() {
	initDBSetting()
}

func main() {
	app, err := NewTelegramBot(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", app.bot.Self.UserName)

	http.HandleFunc("/sns", app.snsHandler)

	updates := app.bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "type /register or /unregister."
		case "register":
			cmd := strings.Split(update.Message.Text, " ")
			if len(cmd) == 0 {
				msg.Text = "I don't know that command"
			} else if strings.Compare(cmd[0], "/register") != 0 {
				msg.Text = "I don't know that command"
			} else {
				token := TokenTbl{
					Name:   cmd[1],
					ChatID: update.Message.Chat.ID,
				}
				if InsertToken(&token) {
					msg.Text = "register ok !!"
				} else {
					msg.Text = "register fail !!"
				}
			}
		case "unregister":
			token := TokenTbl{
				ChatID: update.Message.Chat.ID,
			}
			if DeleteToken(&token) {
				msg.Text = "unregister ok !!"
			} else {
				msg.Text = "unregister fail !!"
			}
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := app.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
