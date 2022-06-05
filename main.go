package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
)

var (
	//todo убрать в структуру env для dependency injection
	DB  *sql.DB
	bot *tgbotapi.BotAPI

	databaseIsActive    = true
	telegramBotToken    string
	openweathermapToken string

	chatovoyNames = map[string]bool{
		shortName: true,
		fullName:  true,
		botName:   true,
		cuteName:  true,
	}
)

const (
	testChatId      = -790845206
	govnosoftChatId = -755317706
	balconyChatId   = -1001416816634

	numberOfKuzyasPictures = 7

	emptyLine = "\n\n"

	greetings                     = "Привет, меня зовут Кузькой, можно Кузенькой. Я маленький ещё, семь веков всего, восьмой пошёл."
	cuteness                      = "😳\U0001F97A😳\U0001F97A😳\U0001F97A"
	showYourselfPhotoErrorMessage = "Ой! Стесняюсь я"
	showYourselfMessage           = "туточки я"
	thankYouResponse              = "Я просто делаю свою работу. Работать буду по совести. За хозяйство не бойся. Конюшня есть?"
	errorMessageDefault           = "Ошибка!"
	skillsIntroduction            = "а вот что я умею:"
	foolMessage                   = "Сам дурак."

	shortName = "чтв"
	fullName  = "чатовой"
	botName   = "@chatovoybot"
	cuteName  = "солнышко заинька"

	specialPlace = "балкон"
)

func main() {
	telegramBotToken = os.Getenv("TELEGRAMBOT_TOKEN")
	openweathermapToken = os.Getenv("OPENWEATHERMAP_TOKEN")
	initBot()
	initDB()

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

		logUpdate(&update)

		reply := ""
		if update.Message.IsCommand() {
			reply = processCommand(&update)
		} else {
			reply = processMessage(&update)
		}

		sendMessage(update.Message.Chat.ID, reply)
	}
}
