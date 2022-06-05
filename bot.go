package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func initBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func sendPhoto(chatID int64, photoBytes []byte, photoName string) {
	photoFileBytes := tgbotapi.FileBytes{
		Name:  photoName,
		Bytes: photoBytes,
	}

	bot.Send(tgbotapi.NewPhoto(chatID, photoFileBytes))
}
