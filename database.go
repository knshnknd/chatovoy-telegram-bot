package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func initDB() {
	databaseHost := os.Getenv("DB_HOST")
	databasePost := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", databaseHost, databasePost, username, password, databaseName)
	fmt.Println(psqlInfo)
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

func requestsCountBySkillAndUser(skill string, userId int64) (int, error) {
	var number int

	if databaseIsActive {
		err := DB.QueryRow("SELECT COUNT(*) FROM requests WHERE skill_name=$1 AND sender_id=$2", skill, userId).Scan(&number)
		if err != nil {
			return 0, err
		}
	}

	return number, nil
}
