package main

import (
  "math/rand"
  "time"
)

func getRandomAnswer() string {

  answers := []string {
		"да.",
		"ну, уж нет. Спасибо этому дому, пойду к другому!",
		"ёщё чего!",
		"не-а!",
		"ой, тошно мне. Нет, конечно.",
		"с большой вероятностью да.",
		"возможно.",
		"вроде, нет.",
		"не знаю, не знаю.",
		"спроси потом ещё раз, мне сейчас не до этого.",
	}
  rand.Seed(time.Now().UnixNano())
  return "Мой ответ: " + answers[rand.Intn(len(answers))]
}
