package main

import (
	"fmt"
	"time"
)

var (
	fullFormat  = "Monday, 2 January 2006 15:04"
	shortFormat = "3:04 PM"
)

func getTime() string {
	timeTbilisi := getTimeForRegion("Asia/Tbilisi", "Tbilisi")
	timeMoscow := getTimeForRegion("Europe/Moscow", "Moscow")
	timeKrasnoyarsk := getTimeForRegion("Asia/Krasnoyarsk", "Krasnoyarsk")
	timeBeijing := getTimeForRegion("Asia/Shanghai", "Shanghai")

	return timeTbilisi + "\n" + timeMoscow + "\n" + timeKrasnoyarsk + "\n" + timeBeijing
}

func getTimeForRegion(locationName, name string) string {
	location, _ := time.LoadLocation(locationName)
	currentLocationTime := time.Now().In(location)

	localTimeMessage := fmt.Sprintf("%s: %s or %s",
		name, currentLocationTime.Format(fullFormat), currentLocationTime.Format(shortFormat))

	return localTimeMessage
}
