package utils

import (
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

type ssdt struct {
	Tyler    string
	Madina   string
	Sreynit  string
	Finn     string
	Paw      string
	Ala      string
	Brian    string
	Anderson string
	Fleur    string
	Karina   string
}

func ReadFile() (string, string) {
	ssdt := map[string]string{
		"Tyler":    "U033S17C8EB",
		"Madina":   "U03GCRZNN6B",
		"Ala":      "U03ENMK3WQ6",
		"Sreynit":  "U013G71SCSC",
		"Finn":     "U03GCRZH4TH",
		"Fleur":    "U03GTDAFWDR",
		"Paw":      "U03GCRZKF9V",
		"Karina":   "U024BST1G80",
		"Brian":    "UN0PANR08",
		"Anderson": "U03HEHMV38F",
	}
	filename := "reviewSchedule"
	f, err := excelize.OpenFile(filename + ".xlsx")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		partner1 := row[1]
		partner1_id := ssdt[partner1]
		partner2 := ""
		partner2_id := ""
		if len(row) > 2 {
			partner2 = row[2]
		}

		if partner2 != "" {
			partner2_id = ssdt[partner2]
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
