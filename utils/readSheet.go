package utils

import (
	"log"

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

func main() {
	ssdt := ssdt{
		Tyler:    "U033S17C8EB",
		Madina:   "U03GCRZNN6B",
		Ala:      "U03ENMK3WQ6",
		Sreynit:  "U013G71SCSC",
		Finn:     "U03GCRZH4TH",
		Fleur:    "U03GTDAFWDR",
		Paw:      "U03GCRZKF9V",
		Karina:   "U024BST1G80",
		Brian:    "UN0PANR08",
		Anderson: "U03HEHMV38F",
	}

	f, err := excelize.OpenFile("reviewSchedule.xlsx")

	if err != nil {
		log.Fatal(err)
	}

}
