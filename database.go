package main

import (
	"database/sql"
	"log"
)

func initDB() {
	psqlInfo := "host=localhost port=54320 user=sandbox " +
		"password=sandbox dbname=sandbox sslmode=disable"

	var err error
	DB, err = sql.Open("pgx", psqlInfo)

	if err != nil {
		log.Println(err)
		turnOffDbFeatures()
	}

	err = DB.Ping()
	if err != nil {
		log.Println(err)
		turnOffDbFeatures()
	}
}

func turnOffDbFeatures() {
	databaseIsActive = false

	//turn off database dependent skills
	existingSkills[bonusesSkill] = false
}

func requestsCountBySkillAndUser(skill string, userId int64) int {
	var number int

	if databaseIsActive {
		err := DB.QueryRow("SELECT COUNT(*) FROM requests WHERE skill_name=$1 AND sender_id=$2", skill, userId).Scan(&number)
		if err != nil {
			log.Fatal(err)
		}
	}

	return number
}
