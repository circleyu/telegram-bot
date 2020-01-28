package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const subConfrmType = "SubscriptionConfirmation"
const notificationType = "Notification"

// snsHandler function for http server
func (app *TelegramBot) snsHandler(w http.ResponseWriter, r *http.Request) {
	chat := getValue("chat", r)
	if chat == "" {
		log.Printf("Unable to Parse url Query")
		return
	}
	token := SelectToken(chat)

	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to Parse Body")
	}
	log.Printf(string(body))
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Printf("Unable to Unmarshal request")
	}

	data := f.(map[string]interface{})
	log.Println(data["Type"].(string))

	if data["Type"].(string) == subConfrmType {
		subcribeURL := data["SubscribeURL"].(string)
		go confirmSubscription(subcribeURL)
	} else if data["Type"].(string) == notificationType {
		var m interface{}
		err = json.Unmarshal([]byte(data["Message"].(string)), &m)
		if err != nil {
			log.Printf("Unable to Unmarshal Message")
		} else {
			message := m.(map[string]interface{})
			log.Printf("Push message to %d: %s", token.ChatID, message["AlarmName"].(string))
			msg := tgbotapi.NewMessage(token.ChatID, message["AlarmName"].(string))
			if _, err := app.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func confirmSubscription(subcribeURL string) {
	response, err := http.Get(subcribeURL)
	if err != nil {
		log.Printf("Unbale to confirm subscriptions")
	} else {
		log.Printf("Subscription Confirmed sucessfully. %d", response.StatusCode)
	}
}

func getValue(key string, r *http.Request) (data string) {
	vars := r.URL.Query()
	datas, ok := vars[key]
	if !ok {
		data = ""
	} else {
		data = datas[0]
	}
	return
}
