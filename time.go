package main

import (
  "time"
)

func getTime() string {
  today := time.Now()
  timeNow :=  "Tbilisi: " + today.Format("Monday, 2 January 2006 15:04") + " or " + today.Format("3:04 PM")

  mos, _ := time.LoadLocation("Europe/Moscow")

  tMoscow := time.Now().In(mos)
  timeMoscow :=  "Moscow: " + tMoscow.Format("Monday, 2 January 2006 15:04") + " or " + tMoscow.Format("3:04 PM")

  krsk, _ := time.LoadLocation("Asia/Krasnoyarsk")

  tKrsk := time.Now().In(krsk)
  timeKrsk :=  "Krasnoyarsk: " + tKrsk.Format("Monday, 2 January 2006 15:04") + " or " + tKrsk.Format("3:04 PM")

  beijing, _ := time.LoadLocation("Asia/Shanghai")

  tBejing := time.Now().In(beijing)
  timeBeijing :=  "Shanghai: " + tBejing.Format("Monday, 2 January 2006 15:04") + " or " + tBejing.Format("3:04 PM")
  str := timeNow + "\n" + timeMoscow + "\n" + timeKrsk + "\n" + timeBeijing

  return str
}
