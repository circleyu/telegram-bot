package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	initDBSetting()
}

func main() {
	initDB()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/")
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
					Name:  cmd[1],
					Token: strconv.FormatInt(update.Message.Chat.ID, 10),
				}
				if InsertToken(&token) {
					msg.Text = "register ok !!"
				} else {
					msg.Text = "register fail !!"
				}
			}
		case "unregister":
			token := TokenTbl{
				Token: strconv.FormatInt(update.Message.Chat.ID, 10),
			}
			if DeleteToken(&token) {
				msg.Text = "unregister ok !!"
			} else {
				msg.Text = "unregister fail !!"
			}
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
