package main

import (
	"github.com/Syfaro/telegram-bot-api"
	owm "github.com/briandowns/openweathermap"

	"log"
	"os"
	"strconv"
	"strings"
	"flag"
)

var (
	// глобальная переменная, в которой храним токен
	telegramBotToken string
	chatID int64
)

// Open Weather Map API-key
var apiKey = os.Getenv("API_WEATHER_KEY")

func init() {
	// меняем BOT_TOKEN на токен бота от BotFather, в строке принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "BOT_TOKEN", "Telegram Bot Token")
	flag.Parse()

	// без флага не запускаем
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}

func main() {
	// используя токен, создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {log.Panic(err)}

	// пишем об этом в консоль
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг, создаем канал, в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update, вычитываем их и обрабатываем
	for update := range updates {
		// сообщение в Lower Case
		messageFromUser := strings.ToLower(update.Message.Text)
		chatID := strconv.FormatInt(int64(update.Message.Chat.ID), 10)

		var reply string

		// универсальный ответ на любое сообщение
		switch messageFromUser {
		case "дурак":
			reply = "Сам дурак."
		case "айди":
			reply = strconv.FormatInt(int64(update.Message.Chat.ID), 10)
		default:
			reply = ""
		}

		if update.Message == nil {continue}

		// разбирает сообщение на слова, реагируем на первое слово
		splitTextFromMessage := strings.Split(messageFromUser, " ")
		switch splitTextFromMessage[0] {
		case "сколько":
			// считаем слова без слова "сколько"
			reply = "Количество слов в этом сообщении без слова «сколько»: " + strconv.Itoa(len(splitTextFromMessage) - 1)
		case "погода":
			w, err := owm.NewCurrent("F", "ru", apiKey)
			if err != nil {log.Fatalln(err)}

			// пока что для примера - Moscow, а вообще второе слово в сообщении - splitTextFromMessage[1]
			w.CurrentByName("Moscow")
			// строчка не ребаотает :С
			// reply = w
		default:
			reply = ""
		}

		// свитч на обработку комманд, комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет."
		case "getChatID":
			reply = chatID
		}

		// логируем, от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// создаем ответное сообщение и отправляем
		msg := tgbotapi.NewMessage(chatID, reply)
		bot.Send(msg)
	}
}
