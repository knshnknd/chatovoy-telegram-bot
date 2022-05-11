package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func prepareMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) PreparedMessage {
	message := update.Message.Text
	messageLowercase := strings.ToLower(message)
	chatID := update.Message.Chat.ID
	splitTextFromMessage := strings.Split(messageLowercase, " ")
	isReplyForBotMessage := isReplyForBot(update, bot)

	botMention := stringContainsInMap(splitTextFromMessage, chatovoyNames, 0)
	botMentionLength := len(strings.Fields(botMention))

	skillName := stringContainsInMap(splitTextFromMessage, existingSkills, botMentionLength)
	skillNameLength := len(strings.Fields(skillName))

	parameter := parseParam(splitTextFromMessage, skillNameLength+botMentionLength)

	return PreparedMessage{
		chatID,
		message,
		messageLowercase,
		splitTextFromMessage,
		isReplyForBotMessage,
		botMention,
		skillName,
		parameter,
	}
}

func isReplyForBot(update tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	reply := update.Message.ReplyToMessage
	if reply == nil || reply.From == nil {
		return false
	}

	return reply.From.ID == bot.Self.ID
}

func stringContainsInMap(splitTextFromMessage []string, arrayOfResults map[string]bool, startIndex int) string {
	result := ""
	currentString := ""
	if startIndex < len(splitTextFromMessage) {
		for i := startIndex; i < len(splitTextFromMessage); i++ {
			if currentString != "" {
				currentString += " "
			}

			currentString += splitTextFromMessage[i]
			if arrayOfResults[currentString] {
				result = currentString
			}
		}
	}

	return result
}

func parseParam(splitTextFromMessage []string, startIndex int) string {
	parameter := ""
	if startIndex < len(splitTextFromMessage) {
		for i := startIndex; i < len(splitTextFromMessage); i++ {
			if parameter != "" {
				parameter += " "
			}
			parameter += splitTextFromMessage[i]
		}
	}

	return parameter
}

type Skill struct {
	name        string
	description string
}

type PreparedMessage struct {
	fromChat             int64
	originalMessage      string
	lowercaseMessage     string
	splitMessage         []string
	isReplyForBotMessage bool

	botMention     string
	skillName      string
	skillParameter string
}
