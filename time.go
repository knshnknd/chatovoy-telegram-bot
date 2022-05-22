package main

import (
	"fmt"
	"time"
)

func getTime() string {
	timeTbilisi := getTimeForRegion("Asia/Tbilisi", "Tbilisi")
	timeMoscow := getTimeForRegion("Asia/Nicosia", "Limassol")
	timeKrasnoyarsk := getTimeForRegion("Asia/Krasnoyarsk", "Krasnoyarsk")
	timeBeijing := getTimeForRegion("Asia/Shanghai", "Shanghai")

	return timeMoscow + "\n" + timeTbilisi + "\n" + timeKrasnoyarsk + "\n" + timeBeijing
}

func getTimeForRegion(locationName, name string) string {
	location, _ := time.LoadLocation(locationName)
	currentTime := time.Now().In(location)

	return fmt.Sprintf("%s: %s or %s", name,
		currentTime.Format("Monday, 2 January 2006 15:04"), currentTime.Format("3:04 PM"))
}
