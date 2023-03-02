package controllers

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlashCommand struct {
	EventHandler *socketmode.SocketmodeHandler
}

func NewSlashCommand(eventhandler *socketmode.SocketmodeHandler) SlashCommand {
	c := SlashCommand{
		EventHandler: eventhandler,
	}

	// Register callback for the command /rocket
	c.EventHandler.HandleSlashCommand(
		"/prs",
		c.getPullRequestData,
	)

	return c

}

func (c SlashCommand) getPullRequestData(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socketmode.Event into a Slash Command
	command, ok := evt.Data.(slack.SlashCommand)

	if !ok {
		log.Printf("ERROR converting event to Slash Command: %v", ok)
	}

	// Make sure to respond to the server to avoid an error
	clt.Ack(*evt.Request)

	// parse the command line
	fmt.Println(command)

	// Post ephemeral message
	_, _, err := clt.PostMessage(
		command.ChannelID,
		slack.MsgOptionResponseURL(command.ResponseURL, slack.ResponseTypeEphemeral),
	)

	// Handle errors
	if err != nil {
		log.Printf("ERROR while sending message for /rocket: %v", err)
	}

}
