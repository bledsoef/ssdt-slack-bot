package controllers

import (
	"log"
	"p01-individual-project-bledsoef/utils"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SocketMode struct {
	EventHandler *socketmode.SocketmodeHandler
}

func NewSlashCommand(eventhandler *socketmode.SocketmodeHandler) SocketMode {
	c := SocketMode{
		EventHandler: eventhandler,
	}

	// Register callback for the command /prs
	c.EventHandler.HandleSlashCommand(
		"/prs",
		c.getPullRequestData,
	)

	// Register callback for the command /cleaning
	c.EventHandler.HandleSlashCommand(
		"/cleaning",
		c.getCleaningData,
	)

	return c

}

func (c SocketMode) getPullRequestData(evt *socketmode.Event, clt *socketmode.Client) {
	// cast our socketmode.Event into a Slash Command
	command, ok := evt.Data.(slack.SlashCommand)

	if !ok {
		log.Printf("ERROR converting event to Slash Command: %v", ok)
	}

	// respond to the server to avoid an error
	clt.Ack(*evt.Request)

	attachment := slack.Attachment{}

	attachment.Text = "Gathering PR data. This may take a moment!"

	clt.PostMessage(
		command.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionResponseURL(command.ResponseURL, slack.ResponseTypeEphemeral),
	)

	message := utils.GetOutstandingPRs()

	attachment.Text = message

	// Post ephemeral message
	_, _, err := clt.PostMessage(
		command.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionResponseURL(command.ResponseURL, slack.ResponseTypeEphemeral),
	)

	// Handle errors
	if err != nil {
		log.Printf("ERROR while sending message for /prs: %v", err)
	}

}

func (c SocketMode) getCleaningData(evt *socketmode.Event, clt *socketmode.Client) {
	// cast our socketmode.Event into a Slash Command
	command, ok := evt.Data.(slack.SlashCommand)

	if !ok {
		log.Printf("ERROR converting event to Slash Command: %v", ok)
	}

	// respond to the server to avoid an error
	clt.Ack(*evt.Request)

	attachment := slack.Attachment{}

	// get the message to send.
	message := utils.GetChoreList()

	attachment.Text = message

	// Post ephemeral message
	_, _, err := clt.PostMessage(
		command.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionResponseURL(command.ResponseURL, slack.ResponseTypeEphemeral),
	)

	// Handle errors
	if err != nil {
		log.Printf("ERROR while sending message for /cleaning: %v", err)
	}
}
