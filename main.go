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
	"strconv"
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

		//объявляем переменные которые понадобятся нам для обработки сообщения
		var reply string = "Я ничего не понял."
		message := update.Message.Text
		messageLowercase := strings.ToLower(message)
		chatID := update.Message.Chat.ID
		splitTextFromMessage := strings.Split(messageLowercase, " ")
		command := update.Message.Command()

		// логируем, от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, message)

		switch splitTextFromMessage[0] {
		case "покажись":
			reply = showYourself(bot, chatID)
		case "сколько":
			// считаем слова без слова "сколько"
			reply = wordsCount(splitTextFromMessage)
		case "погода":
			if len(splitTextFromMessage) > 2 {
				reply = "Больше двух слов не пиши, когда погоду хочешь узнать!"
			} else {
				reply = requestWeather(splitTextFromMessage)
			}
		case "дурак":
			reply = "Сам дурак."
		case "айди":
			reply = strconv.FormatInt(chatID, 10)
		case "спасибо":
			reply = "Я просто делаю свою работу. Работать буду по совести. За хозяйство не бойся. Конюшня есть?"
		}

		// свитч на обработку комманд, комманда - сообщение, начинающееся с "/"
		switch command {
		case "start":
			reply = "Привет, меня зовут Кузькой, можно Кузенькой. Я маленький ещё, семь веков всего, восьмой пошёл."
		case "getChatID":
			reply = strconv.FormatInt(chatID, 10)
		case "currency":
			reply = getCurrency()
		case "time":
			reply = getTime()
		case "random":
			reply = getRandomAnswer()
		}

		// создаем ответное сообщение и отправляем
		msg := tgbotapi.NewMessage(chatID, reply)
		bot.Send(msg)
	}
}

func wordsCount(splitTextFromMessage []string) string {
	return "Количество слов в этом сообщении без слова «сколько»: " + strconv.Itoa(len(splitTextFromMessage)-1)
}

func requestWeather(splitTextFromMessage []string) string {
	w, err := owm.NewCurrent("C", "ru", openweathermapToken)
	if err != nil {
		log.Fatalln(err)
	}

	if splitTextFromMessage[1] == "балкон" {
		return "На балконе как всегда тепло и уютно."
	} else {
	err = w.CurrentByName(splitTextFromMessage[1])

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

	err = f.DailyByName(splitTextFromMessage[1], 5)

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
