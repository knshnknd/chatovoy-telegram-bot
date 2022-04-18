package main

import (
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	owm "github.com/briandowns/openweathermap"
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
)

var (
	// глобальная переменная, в которой храним токен
	telegramBotToken string
	openweathermapToken string
)

// Open Weather Map API-key
// var apiKey = os.Getenv("a3fb0c63cfb5b617e03f3e7d38b753c1")

func init() {
	// меняем BOT_TOKEN на токен бота от BotFather, в строке принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "5177088641:AAEx4sreMj1o7ZxUqDF3bmD0rIg0X-a6l-U", "Telegram Bot Token")
	flag.StringVar(&openweathermapToken, "openweathermapToken", "a3fb0c63cfb5b617e03f3e7d38b753c1", "OpenWeatherMap Token")
	flag.Parse()

	// без флага не запускаем
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}

func main() {
	// используя токен, создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// пишем об этом в консоль
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг, создаем канал, в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update, вычитываем их и обрабатываем
	for update := range updates {

		if update.Message == nil {
			continue
		}

		//объявляем переменные которые понадобятся нам для обработки сообщения
		var reply string = "я ничего не понял"
		message := update.Message.Text
		messageLowercase := strings.ToLower(message)
		chatID := update.Message.Chat.ID
		splitTextFromMessage := strings.Split(messageLowercase, " ")
		command := update.Message.Command()

		// логируем, от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, message)

		switch splitTextFromMessage[0] {
		case "сколько":
			// считаем слова без слова "сколько"
			reply = "Количество слов в этом сообщении без слова «сколько»: " + strconv.Itoa(len(splitTextFromMessage)-1)
		case "погода":
			w, err := owm.NewCurrent("C", "ru", openweathermapToken)
			if err != nil {
			log.Fatalln(err)
			}

			w.CurrentByName(splitTextFromMessage[1])

			reply = "Погода:" + w.Name + ": " + fmt.Sprintf("%.1f", w.Main.Temp) + "°C, " + w.Weather[0].Description + ", влажность: " + strconv.Itoa(w.Main.Humidity) + "%"
		case "дурак":
			reply = "Сам дурак."
		case "айди":
			reply = strconv.FormatInt(chatID, 10)
		}

		// свитч на обработку комманд, комманда - сообщение, начинающееся с "/"
		switch command {
		case "start":
			reply = "Привет."
		case "getChatID":
			reply = strconv.FormatInt(chatID, 10)
		}

		// создаем ответное сообщение и отправляем
		msg := tgbotapi.NewMessage(chatID, reply)
		bot.Send(msg)
	}
}
