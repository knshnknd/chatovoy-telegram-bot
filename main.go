package main

import (
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
)

var (
	// глобальная переменная, в которой храним токен
	telegramBotToken    string
	openweathermapToken string

	skills = []Skill{
		{name: "расскажись", description: "ну вы уже поняли как оно работает"},
		{name: "покажись", description: "явлюсь к вам во всей своей красе"},
		{name: "ответь", description: "с вас вопрос с меня ответ"},
		{name: "погода", description: "выгляну в окно за вас"},
		{name: "дурак", description: "даже не думай"},
		{name: "спасибо", description: "вежливость у нас в почёте"},
	}

	existingSkills = map[string]bool{
		"покажись": true,
		"ответь":   true,
		"погода":   true,
		"дурак":    true,
		"спасибо":  true,
	}

	chatovoyNames = map[string]bool{
		"чтв":              true,
		"чатовой":          true,
		"@chatovoybot":     true,
		"солнышко заинька": true,
	}
)

const (
	testChatId      = -790845206
	govnosoftChatId = -755317706
	balconyChatId   = -1001416816634

	numberOfKuzyasPictures = 7

	emptyLine = "\n\n"
	greetings = "Привет, меня зовут Кузькой, можно Кузенькой. Я маленький ещё, семь веков всего, восьмой пошёл."
	cuteness  = "😳\U0001F97A😳\U0001F97A😳\U0001F97A"
)

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

		sendMessage(bot, update.Message.Chat.ID, reply)
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
		reply = greetings
	case "currency":
		reply = getCurrency()
	case "time":
		reply = getTime()
	}
	return reply
}

func processMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	message := prepareMessage(update, bot)
	reply := ""

	if isMessageForBot(message) {
		switch message.skillName {
		case "расскажись":
			reply = introduceYourself()
		case "покажись":
			reply = showYourself(bot, message.fromChat)
		case "ответь":
			reply = getRandomAnswer()
		case "погода":
			reply = showWeather(message.skillParameter)
		case "дурак":
			reply = "Сам дурак."
		case "спасибо":
			reply = "Я просто делаю свою работу. Работать буду по совести. За хозяйство не бойся. Конюшня есть?"
		}
	}

	if message.botMention == "солнышко заинька" {
		reply += " " + cuteness
	}

	return reply
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func introduceYourself() string {
	skillsIntroduction := "а вот что я умею:"

	skillsText := ""

	for _, elem := range skills {
		skillsText += fmt.Sprintf("%s -> %s\n", elem.name, elem.description)
	}

	return greetings + emptyLine + skillsIntroduction + emptyLine + skillsText
}

func showWeather(place string) string {
	return requestWeatherByPlace(place)
}

func showYourself(bot *tgbotapi.BotAPI, chatID int64) string {
	reply := "туточки я"

	photoName := generatePhotoName()
	photoBytes, err := ioutil.ReadFile(makePhotoPath(photoName))

	if err != nil {
		return "Ой! Стесняюсь я"
	} else {
		sendPhoto(bot, chatID, photoBytes, photoName)
		return reply
	}
}

func isMessageForBot(message PreparedMessage) bool {
	return chatovoyNames[message.botMention] || existingSkills[message.skillName] && message.isReplyForBotMessage
}
