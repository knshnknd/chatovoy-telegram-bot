package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"log"
)

func requestWeatherByPlace(place string) string {
	w, err := owm.NewCurrent("C", "ru", openweathermapToken)
	if err != nil {
		log.Fatalln(err)
	}

	if place == "балкон" {
		return "На балконе как всегда тепло и уютно."
	} else {
		err = w.CurrentByName(place)

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

		err = f.DailyByName(place, 5)

		if err != nil {
			return "Ошибка!"
		}

		return currentWeather
	}
}
