package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type PreparedMessage struct {
	senderId             int64
	fromChat             int64
	originalMessage      string
	isReplyForBotMessage bool

	botMention     string
	skillName      string
	skillParameter string
}

func prepareMessage(update *tgbotapi.Update) PreparedMessage {
	userId := update.Message.From.ID
	message := update.Message.Text
	chatID := update.Message.Chat.ID
	isReplyForBotMessage := isReplyForBot(update)

	lowercaseMessage := strings.ToLower(message)
	splitTextFromMessage := strings.Split(lowercaseMessage, " ")

	botMention := stringContainsInMap(splitTextFromMessage, chatovoyNames, 0)
	botMentionLength := len(strings.Fields(botMention))

	skillName := stringContainsInMap(splitTextFromMessage, existingSkills, botMentionLength)
	skillNameLength := len(strings.Fields(skillName))

	parameter := parseParam(splitTextFromMessage, skillNameLength+botMentionLength)

	return PreparedMessage{
		userId,
		chatID,
		message,
		isReplyForBotMessage,
		botMention,
		skillName,
		parameter,
	}
}

func isReplyForBot(update *tgbotapi.Update) bool {
	reply := update.Message.ReplyToMessage
	if reply == nil || reply.From == nil {
		return false
	}

	return reply.From.ID == bot.Self.ID
}

func isMessageForBot(message PreparedMessage) bool {
	return chatovoyNames[message.botMention] || existingSkills[message.skillName] && message.isReplyForBotMessage
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
