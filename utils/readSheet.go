package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

func ReadReviewSchedule() (string, string) {
	// open excel file
	filename := "ssdtSchedule"
	f, err := excelize.OpenFile(filename + ".xlsx")

	// handle potential erros
	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("Review")

	// handle potential errors
	if err != nil {
		log.Fatal(err)
	}

	// iterate through each row in the sheet and return the schedule to see who the two current reviewers are.
	for _, row := range rows {
		partner1 := row[1]
		partner1_id := GetUserID(partner1)
		partner2 := ""
		partner2_id := ""
		if len(row) > 2 {
			partner2 = row[2]
		}

		if partner2 != "" {
			partner2_id = GetUserID(partner2)
		}

		// if it is the correct date column
		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			log.Fatal(err)
		}

		if time.Since(date) < time.Hour*168 {
			return partner1_id, partner2_id
		}

	}
	return "", ""
}

func ReadChoreSchedule() [][]string {
	choreList := [][]string{}

	// open the excel file
	filename := "ssdtSchedule"
	f, err := excelize.OpenFile(filename + ".xlsx")

	// handle potential error
	if err != nil {
		log.Fatal(err)
	}
	index := 0
	rows, _ := f.GetRows("Chores")

	// iterate through the first row in the sheet to grab the row index of the current date
	for i, row := range rows[0][1:] {
		date, err := time.Parse("2006-01-02", row)
		if err != nil {
			log.Fatal(err)
		}
		if time.Since(date) < time.Hour*24 {
			index = i
		}

	}

	// catch potential errors
	cols, err := f.GetCols("Chores")
	if err != nil {
		log.Fatal(err)
	}
	// iterate through the columns
	for i, col := range cols {
		// if the column is the current date
		if i == index {
			// get the values of the column except for the date
			people := col[1:]
			// iterate through all the people
			for j, person := range people {
				userChore := cols[0][j+1]
				user := person
				babyList := []string{user, userChore}
				choreList = append(choreList, babyList)
			}
		}
	}
	return choreList
}

func GetUserID(name string) string {
	// initialize a map that maps both the first name and github usernames to the Slack User IDs respectively.
	ssdtMap := map[string]string{
		"Tyler":      "U033S17C8EB",
		"tylerpar99": "U033S17C8EB",

		"Madina":      "U03GCRZNN6B",
		"solijonovam": "U03GCRZNN6B",

		"Ala":    "U03ENMK3WQ6",
		"qasema": "U03ENMK3WQ6",

		"Sreynit":   "U013G71SCSC",
		"sreynit02": "U013G71SCSC",

		"Finn":     "U03GCRZH4TH",
		"bledsoef": "U03GCRZH4TH",

		"Fleur":      "U03GTDAFWDR",
		"gahimbaref": "U03GTDAFWDR",

		"Paw":    "U03GCRZKF9V",
		"thawpt": "U03GCRZKF9V",

		"Karina":            "U024BST1G80",
		"Karina-Agliullova": "U024BST1G80",

		"Brian":       "UN0PANR08",
		"BrianRamsay": "UN0PANR08",

		"Anderson":         "U03HEHMV38F",
		"Andersonstettner": "U03HEHMV38F",
	}

	return ssdtMap[name]
}

func GetChoreList() string {
	message := "*This week's cleaning assignments are:* \n"
	choreList := ReadChoreSchedule()
	for _, chore := range choreList {
		message += fmt.Sprintf("*%s:* %s <@%s> \n", chore[0], chore[1], GetUserID(chore[0]))
	}
	return message
}
