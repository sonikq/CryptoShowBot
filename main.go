package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type binanceResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("5620373699:AAEBR0v1khgLCeSP9YTFWFplW7PA9AxixLk")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msgArr := strings.Split(update.Message.Text, " ")

			switch msgArr[0] {
			case "SHOW":
				ans := "Курс валюты(USDT):\n"

				coinPrice, err := getPrice(msgArr[1])
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
				}
				ans += fmt.Sprintf("[%.2f]", coinPrice)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная валюта")
				bot.Send(msg)
			}
		}
	}
}
func getPrice(coin string) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var jsonResp binanceResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	price = jsonResp.Price
	if err != nil {
		return
	}

	if jsonResp.Code != 0 {
		err = errors.New("Указана неверная валюта")
		return
	}

	return
}
