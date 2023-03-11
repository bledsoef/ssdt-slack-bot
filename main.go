package main

import (
	"log"
	"os"

	"p01-individual-project-bledsoef/controllers"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	// initialize slack instance
	slackClient := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	// initialize socket mode instance
	socket := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	socketmodeHandler := socketmode.NewSocketmodeHandler(socket)

	controllers.NewSlashCommand(socketmodeHandler)

	// controllers.NewEvent(socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}
