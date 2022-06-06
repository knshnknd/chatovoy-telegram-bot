package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Skill struct {
	name        string
	description string
}

var (
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
)

const (
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
)

func howManyThankYou(userId int64) string {
	if databaseIsActive {
		counter, err := requestsCountBySkillAndUser(thankYouSkill, userId)
		if err != nil {
			return errorMessageDefault
		}
		return fmt.Sprintf("Спасиб на вашем счету: %d", counter)
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

func showYourself(chatID int64) string {
	reply := showYourselfMessage

	photoName := generatePhotoName()
	photoBytes, err := ioutil.ReadFile(makePhotoPath(photoName))

	if err != nil {
		return showYourselfPhotoErrorMessage
	} else {
		sendPhoto(chatID, photoBytes, photoName)
		return reply
	}
}

func requestWeatherByPlace(place string) string {
	if place == specialPlace {
		return balconyWeather()
	} else {
		weatherByPlace, err := weatherByPlace(place)
		if err != nil {
			return errorMessageDefault
		}
		return oneCityResponse(weatherByPlace)
	}
}

func getTime() string {
	timeTbilisi := getTimeForRegion("Asia/Tbilisi", "Tbilisi")
	timeMoscow := getTimeForRegion("Asia/Nicosia", "Limassol")
	timeKrasnoyarsk := getTimeForRegion("Asia/Krasnoyarsk", "Krasnoyarsk")
	timeBeijing := getTimeForRegion("Asia/Shanghai", "Shanghai")

	return timeMoscow + "\n" + timeTbilisi + "\n" + timeKrasnoyarsk + "\n" + timeBeijing
}

func getRandomAnswer() string {
	rand.Seed(time.Now().UnixNano())
	return "Мой ответ: " + answers[rand.Intn(len(answers))]
}

func getCurrency() string {
	resp, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	xmlCurrenciesData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var env Envelope
	err = xml.Unmarshal(xmlCurrenciesData, &env)
	if err != nil {
		log.Fatal(err)
	}

	var reply string
	for _, v := range env.Cube[0].Rates {
		if v.Currency == "USD" || v.Currency == "JPY" || v.Currency == "GBP" {
			reply = reply + "1 евро = " + v.Rate + " " + v.Currency + "\n"
		}
	}

	return reply
}
