package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/oauth2"
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

	return c

}

// func NewEvent(eventhandler *socketmode.SocketmodeHandler) SocketMode {
// 	c := SocketMode{
// 		EventHandler: eventhandler,
// 	}

// 	// Register callback for the command /prs
// 	c.EventHandler.HandleEvents(
// 		slackevents.AppMention,
// 		c.appMention,
// 	)

// 	return c

// }

func (c SocketMode) getPullRequestData(evt *socketmode.Event, clt *socketmode.Client) {
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

// func (c SocketMode) appMention(evt *socketmode.Event, clt *socketmode.Client) {
// 	evt_api, _ := evt.Data.(slackevents.EventsAPIEvent)
// 	evt_app_mention, _ := evt_api.InnerEvent.Data.(*slackevents.AppMentionEvent)
// 	fmt.Println("HELLO", evt_app_mention., "Hello")
// }

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
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
