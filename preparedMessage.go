package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func prepareMessage(update tgbotapi.Update) PreparedMessage {
	message := update.Message.Text
	messageLowercase := strings.ToLower(message)
	chatID := update.Message.Chat.ID
	splitTextFromMessage := strings.Split(messageLowercase, " ")

	botMention, endOfBotMention := parseName(splitTextFromMessage)
	skillName, endOfSkillName := parseSkill(splitTextFromMessage, endOfBotMention+1)
	parameter := parseParam(splitTextFromMessage, endOfSkillName+1)

	return PreparedMessage{
		chatID,
		message,
		messageLowercase,
		splitTextFromMessage,
		botMention,
		skillName,
		parameter,
	}
}

func parseName(splitTextFromMessage []string) (string, int) {
	name := ""
	nameIndexEnd := 0
	for index, element := range splitTextFromMessage {
		if name != "" {
			name += " "
		}

		name += element
		if chatovoyNames[name] {
			nameIndexEnd = index
			return name, nameIndexEnd
		}
	}

	return name, nameIndexEnd
}

func parseSkill(splitTextFromMessage []string, startIndex int) (string, int) {
	skill := ""
	nameIndexEnd := startIndex
	if startIndex < len(splitTextFromMessage) {
		for i := startIndex; i < len(splitTextFromMessage); i++ {
			if skill != "" {
				skill += " "
			}

			skill += splitTextFromMessage[i]
			if existingSkills[skill] {
				nameIndexEnd = i
				return skill, nameIndexEnd
			}
		}
	}

	return skill, nameIndexEnd
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
	fromChat         int64
	originalMessage  string
	lowercaseMessage string
	splitMessage     []string

	botMention     string
	skillName      string
	skillParameter string
}
