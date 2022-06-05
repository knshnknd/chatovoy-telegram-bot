package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io/ioutil"
	"log"
	"os"
)

var (
	DB                  *sql.DB
	databaseIsActive    = true
	telegramBotToken    string
	openweathermapToken string

	skills = []Skill{
		{name: introduceSkill, description: "ну вы уже поняли как оно работает"},
		{name: showSkill, description: "явлюсь к вам во всей своей красе"},
		{name: answerSkill, description: "с вас вопрос с меня ответ"},
		{name: weatherSkill, description: "выгляну в окно за вас"},
		{name: youFoolSkill, description: "даже не думай"},
		{name: thankYouSkill, description: "вежливость у нас в почёте"},
		{name: currencyCommand, description: "невнятный курс валют без любимого рублика"},
		{name: timeCommand, description: "текущее время в главных городах мира"},
		{name: bonusesSkill, description: "узнаем насколько ты благодарный"},
	}

	existingSkills = map[string]bool{
		introduceSkill: true,
		showSkill:      true,
		answerSkill:    true,
		weatherSkill:   true,
		youFoolSkill:   true,
		thankYouSkill:  true,
		bonusesSkill:   true,
	}

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

	introduceSkill = "расскажись"
	showSkill      = "покажись"
	answerSkill    = "ответь"
	weatherSkill   = "погода"
	youFoolSkill   = "дурак"
	thankYouSkill  = "спасибо"
	bonusesSkill   = "сколько у меня спасиб"

	startCommand    = "start"
	currencyCommand = "currency"
	timeCommand     = "time"

	shortName = "чтв"
	fullName  = "чатовой"
	botName   = "@chatovoybot"
	cuteName  = "солнышко заинька"

	specialPlace = "балкон"
)

func initDB() error {
	psqlInfo := "host=localhost port=54320 user=sandbox " +
		"password=sandbox dbname=sandbox sslmode=disable"

	var err error
	DB, err = sql.Open("pgx", psqlInfo)

	if err != nil {
		log.Println(err)
		turnOffDbFeatures()
	}

	return DB.Ping()
}

func turnOffDbFeatures() {
	databaseIsActive = false

	//turn off database dependent skills
	existingSkills[bonusesSkill] = false
}

func main() {
	err := initDB()
	if err != nil {
		log.Println(err)
		turnOffDbFeatures()
	}

	telegramBotToken = os.Getenv("TELEGRAMBOT_TOKEN")
	openweathermapToken = os.Getenv("OPENWEATHERMAP_TOKEN")
	fmt.Println(telegramBotToken)
	fmt.Println(openweathermapToken)

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
	case startCommand:
		reply = greetings
	case currencyCommand:
		reply = getCurrency()
	case timeCommand:
		reply = getTime()
	}
	return reply
}

func processMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	message := prepareMessage(update, bot)
	reply := ""

	if isMessageForBot(message) {
		logMessage(&message)
		switch message.skillName {
		case introduceSkill:
			reply = introduceYourself()
		case showSkill:
			reply = showYourself(bot, message.fromChat)
		case answerSkill:
			reply = getRandomAnswer()
		case weatherSkill:
			reply = showWeather(message.skillParameter)
		case youFoolSkill:
			reply = foolMessage
		case thankYouSkill:
			reply = thankYouResponse
		case bonusesSkill:
			reply = howManyThankYou(message.senderId)
		}
	}

	if message.botMention == cuteName {
		reply += emptyLine + cuteness
	}

	return reply
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func howManyThankYou(userId int64) string {
	if databaseIsActive {
		return fmt.Sprintf("Спасиб на вашем счету: %d", skillCount(thankYouSkill, userId))
	}

	return ""
}

func introduceYourself() string {
	skillsText := ""

	for _, elem := range skills {
		if existingSkills[elem.name] {
			skillsText += fmt.Sprintf("%s -> %s\n", elem.name, elem.description)
		}
	}

	return greetings + emptyLine + skillsIntroduction + emptyLine + skillsText
}

func showWeather(place string) string {
	return requestWeatherByPlace(place)
}

func showYourself(bot *tgbotapi.BotAPI, chatID int64) string {
	reply := showYourselfMessage

	photoName := generatePhotoName()
	photoBytes, err := ioutil.ReadFile(makePhotoPath(photoName))

	if err != nil {
		return showYourselfPhotoErrorMessage
	} else {
		sendPhoto(bot, chatID, photoBytes, photoName)
		return reply
	}
}

func logMessage(message *PreparedMessage) {
	if databaseIsActive {
		query :=
			`INSERT INTO requests(
                     sender_id, 
                     from_chat, 
                     original_message, 
                     is_reply, 
                     bot_mention, 
                     skill_name, 
                     skill_parameter) 
		 	VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := DB.Exec(query,
			message.senderId,
			message.fromChat,
			message.originalMessage,
			message.isReplyForBotMessage,
			message.botMention,
			message.skillName,
			message.skillParameter,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func skillCount(skill string, userId int64) int {
	var number int

	if databaseIsActive {
		err := DB.QueryRow("SELECT COUNT(*) FROM requests WHERE skill_name=$1 AND sender_id=$2", skill, userId).Scan(&number)
		if err != nil {
			log.Fatal(err)
		}
	}

	return number
}
