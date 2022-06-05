package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"log"
)

func weatherByPlace(place string) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent("C", "ru", openweathermapToken)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.CurrentByName(place)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func balconyWeather() string {
	limassol, _ := weatherByPlace("лимасол")
	rustavi, _ := weatherByPlace("рустави")
	krasnoyarsk, _ := weatherByPlace("красноярск")
	chzenchzou, _ := weatherByPlace("чжэнчжоу")

	return specialPlacesResponse(krasnoyarsk, chzenchzou, rustavi, limassol)
}

func specialPlacesResponse(weathers ...*owm.CurrentWeatherData) string {
	response := "погода в самых важных местах"

	for _, a := range weathers {
		response += emptyLine
		switch a.Name {
		case "Лимасол":
			response += "🏝"
		case "Рустави":
			response += "⛰"
		case "Красноярск":
			response += "🏠"
		case "Чжэнчжоу":
			response += "🏙"
		}
		response += specialPlaceResponse(a)
	}

	return response
}

func oneCityResponse(w *owm.CurrentWeatherData) string {
	return fmt.Sprintf("Погода в городе %s: %.1f °C, %s, влажность: %d%%",
		w.Name, w.Main.Temp, w.Weather[0].Description, w.Main.Humidity)
}

func specialPlaceResponse(w *owm.CurrentWeatherData) string {
	return fmt.Sprintf("Балкон.%s: %.1f °C, %s, влажность: %d%%",
		w.Name, w.Main.Temp, w.Weather[0].Description, w.Main.Humidity)
}

//todo сделать прогноз на 5 дней
func forecastForDays() string {
	f, err := owm.NewForecast("5", "C", "ru", openweathermapToken)
	if err != nil {
		log.Fatalln(err)
	}

	err = f.DailyByName("place", 5)

	if err != nil {
		return errorMessageDefault
	}

	return ""
}
