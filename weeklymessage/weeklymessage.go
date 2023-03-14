package main

import (
	"os"
	"p01-individual-project-bledsoef/utils"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")

	channelID := "C04LA97FWKH"

	api := slack.New(token)
	if os.Args[1] == "pr" {
		api.PostMessage(
			channelID,
			slack.MsgOptionText(utils.GetOutstandingPRs(), false),
		)
	} else if os.Args[1] == "clean" {
		api.PostMessage(
			channelID,
			slack.MsgOptionText(utils.GetChoreList(), false),
		)
	}

}
