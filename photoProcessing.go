package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"math/rand"
)

func generatePhotoName() string {
	return fmt.Sprintf("kuzya%d", rand.Intn(numberOfKuzyasPictures))
}

func sendPhoto(bot *tgbotapi.BotAPI, chatID int64, photoName string) {
	photoBytes, err := ioutil.ReadFile(makePhotoPath(photoName))

	if err != nil {
		panic(err)
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  photoName,
		Bytes: photoBytes,
	}
	bot.Send(tgbotapi.NewPhoto(chatID, photoFileBytes))
}

func makePhotoPath(photoName string) string {
	return fmt.Sprintf("resources/%s.jpg", photoName)
}
