package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func logUpdate(update *tgbotapi.Update) {
	message := update.Message.Text
	userName := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	chatTitle := update.Message.Chat.Title

	log.Printf("[%s] sent message: \"%s\" to chat: \"%s\"[%d]", userName, message, chatTitle, chatID)
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
			log.Println(err)
		}
	}
}

func processCommand(update *tgbotapi.Update) string {
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

func processMessage(update *tgbotapi.Update) string {
	message := prepareMessage(update)
	reply := ""

	if isMessageForBot(message) {
		logMessage(&message)
		switch message.skillName {
		case introduceSkill:
			reply = introduceYourself()
		case showSkill:
			reply = showYourself(message.fromChat)
		case answerSkill:
			reply = getRandomAnswer()
		case weatherSkill:
			reply = requestWeatherByPlace(message.skillParameter)
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
