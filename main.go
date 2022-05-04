package main

import (
	"flag"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

var (
	// глобальная переменная, в которой храним токен
	telegramBotToken    string
	openweathermapToken string
)

// Open Weather Map API-key
// var apiKey = os.Getenv("a3fb0c63cfb5b617e03f3e7d38b753c1")

func init() {
	// меняем BOT_TOKEN на токен бота от BotFather, в строке принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.StringVar(&openweathermapToken, "openweathermapToken", "", "OpenWeatherMap Token")
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

		logUpdate(update)

		reply := ""
		if update.Message.IsCommand() {
			reply = processCommand(update)
		} else {
			reply = processMessage(update, bot)
		}

		sendReplyToUpdate(update, reply, bot)
	}
}

func logUpdate(update tgbotapi.Update) {
	message := update.Message.Text
	userName := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	chatTitle := update.Message.Chat.Title

	log.Printf("[%s] sent message: \"%s\" to chat: \"%s\"[%d]", userName, message, chatTitle, chatID)
}

func processCommand(update tgbotapi.Update) string {
	command := update.Message.Command()
	reply := ""

	switch command {
	case "start":
		reply = "Привет, меня зовут Кузькой, можно Кузенькой. Я маленький ещё, семь веков всего, восьмой пошёл."
	case "currency":
		reply = getCurrency()
	case "time":
		reply = getTime()
	}
	return reply
}

func processMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	message := update.Message.Text
	messageLowercase := strings.ToLower(message)
	chatID := update.Message.Chat.ID
	splitTextFromMessage := strings.Split(messageLowercase, " ")
	reply := ""

	if isMessageForBot(splitTextFromMessage) {
		switch splitTextFromMessage[1] {
		case "покажись":
			reply = showYourself(bot, chatID)
		case "ответь":
			reply = getRandomAnswer()
		case "погода":
			reply = showWeather(splitTextFromMessage)
		case "дурак":
			reply = "Сам дурак."
		case "спасибо":
			reply = "Я просто делаю свою работу. Работать буду по совести. За хозяйство не бойся. Конюшня есть?"
		}
	}

	return reply
}

func sendReplyToUpdate(update tgbotapi.Update, reply string, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}

func showWeather(splitTextFromMessage []string) string {
	reply := ""

	if len(splitTextFromMessage) > 3 {
		reply = "Больше двух слов не пиши, когда погоду хочешь узнать!"
	} else {
		place := splitTextFromMessage[2]
		reply = requestWeatherByPlace(place)
	}
	return reply
}

func isMessageForBot(splitTextFromMessage []string) bool {
	//это тупа имитация HashSet из Java, потому что по-дефолту в Go нет множеств
	chatovoyNames := map[string]bool{
		"чтв":          true,
		"чатовой":      true,
		"@chatovoybot": true,
	}
	return chatovoyNames[splitTextFromMessage[0]]
}

func requestWeatherByPlace(place string) string {
	w, err := owm.NewCurrent("C", "ru", openweathermapToken)
	if err != nil {
		log.Fatalln(err)
	}

	if place == "балкон" {
		return "На балконе как всегда тепло и уютно."
	} else {
		err = w.CurrentByName(place)

		if err != nil {
			return "Ошибка!"
		}

		currentWeather := fmt.Sprintf("Погода в городе %s: %.1f °C, %s, влажность: %d%%",
			w.Name, w.Main.Temp, w.Weather[0].Description, w.Main.Humidity)

		// ПОТОМ СДЕЛАЮ В ОТДЕЛЬНЫЙ ФАЙЛ ВСЮ ПОГОДУ!!!
		f, err := owm.NewForecast("5", "C", "ru", openweathermapToken)
		if err != nil {
			log.Fatalln(err)
		}

		err = f.DailyByName(place, 5)

		if err != nil {
			return "Ошибка!"
		}

		forecastWeather := "Прогноз на 5 дней в разработке..."

		return currentWeather + "\n\n" + forecastWeather
	}
}

func showYourself(bot *tgbotapi.BotAPI, chatID int64) string {
	reply := "туточки я"
	sendPhoto(bot, chatID, generatePhotoName())

	return reply
}

func generatePhotoName() string {
	return fmt.Sprintf("kuzya%d", rand.Intn(7))
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
