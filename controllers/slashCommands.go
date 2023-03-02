package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/oauth2"
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
	repo := command.Text

	attachment := slack.Attachment{}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	// ctx, cancel := context.WithCancel(context.Background())

	message := GetPRs(repo, githubClient, ctx)

	attachment.Text = message

	// Post ephemeral message
	_, _, err := clt.PostMessage(
		command.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionResponseURL(command.ResponseURL, slack.ResponseTypeEphemeral),
	)

	// Handle errors
	if err != nil {
		log.Printf("ERROR while sending message for /rocket: %v", err)
	}

}

func GetPRs(repo string, githubClient *github.Client, ctx context.Context) string {
	message := ""
	if repo == "all" {
		repos := [...]string{"celts", "lsf", "bcsr"}
		// Loop through the list of pull requests and print out the title and URL
		var allprs []*github.PullRequest
		for _, r := range repos {
			prs, _, err := githubClient.PullRequests.List(ctx, "BCStudentSoftwareDevTeam", r, nil)
			allprs = append(allprs, prs...)
			if err != nil {
				fmt.Printf("Error retrieving pull requests: %v\n", err)
			}
		}
		// Loop through the list of pull requests and print out the title and URL
		for i, pr := range allprs {
			message += fmt.Sprintf("%s. %s last updated on %s. View it here: %s \n", strconv.Itoa(i+1), *pr.Title, *pr.UpdatedAt, *pr.HTMLURL)
		}
	} else {
		prs, _, err := githubClient.PullRequests.List(ctx, "BCStudentSoftwareDevTeam", repo, nil)
		if err != nil {
			fmt.Printf("Error retrieving pull requests: %v\n", err)
		}
		// Loop through the list of pull requests and print out the title and URL
		for i, pr := range prs {
			message += fmt.Sprintf("%s. %s last updated on %s. View it here: %s \n", strconv.Itoa(i+1), *pr.Title, *pr.UpdatedAt, *pr.HTMLURL)
		}
	}

	return message
}
