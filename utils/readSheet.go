package utils

import (
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

func ReadReviewSchedule() (string, string) {

	filename := "ssdtSchedule"
	f, err := excelize.OpenFile(filename + ".xlsx")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("Review")
	if err != nil {
		log.Fatal(err)
	}

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
	filename := "ssdtSchedule"
	f, err := excelize.OpenFile(filename + ".xlsx")

	if err != nil {
		log.Fatal(err)
	}
	index := 0
	rows, _ := f.GetRows("Chores")
	for i, row := range rows[0][1:] {
		date, err := time.Parse("2006-01-02", row)
		if err != nil {
			log.Fatal(err)
		}
		if time.Since(date) < time.Hour*24 {
			index = i
		}

	}

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
