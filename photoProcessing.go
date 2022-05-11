package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

func generatePhotoName() string {
	return fmt.Sprintf("kuzya%d", rand.Intn(numberOfKuzyasPictures))
}

func sendPhoto(bot *tgbotapi.BotAPI, chatID int64, photoBytes []byte, photoName string) {
	photoFileBytes := tgbotapi.FileBytes{
		Name:  photoName,
		Bytes: photoBytes,
	}

	bot.Send(tgbotapi.NewPhoto(chatID, photoFileBytes))
}

func makePhotoPath(photoName string) string {
	return fmt.Sprintf("resources/%s.jpg", photoName)
}
