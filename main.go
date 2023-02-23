package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/oauth2"
)

func main() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	slackClient := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socket := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	now := time.Now()
	now.Add(time.Minute * 1)
	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()
	go func(ctx context.Context, slackClient *slack.Client, socket *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {

			attachment := slack.Attachment{}
			attachment.Text = "Hi"
			socket.ScheduleMessage("C04LA97FWKH", now.String(), slack.MsgOptionAttachments(attachment))
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socket.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socket.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := HandleEventMessage(eventsAPI, slackClient)
					if err != nil {
						// TODO: Replace with actual err handeling
						log.Fatal(err)
					}
				}
			}
		}
	}(ctx, slackClient, socket)

	socket.Run()
}

// HandleEventMessage will take an event and handle it properly based on the type of event
func HandleEventMessage(event slackevents.EventsAPIEvent, slackClient *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := HandleAppMentionEventToBot(ev, slackClient)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// HandleAppMentionEventToBot is used to take care of the AppMentionEvent when the bot is mentioned
func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, slackClient *slack.Client) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := slackClient.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	text := strings.ToLower(event.Text)

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	// attachment.Fields = []slack.AttachmentField{
	// 	{
	// 		Title: "Date",
	// 		Value: time.Now().String(),
	// 	}, {
	// 		Title: "Initializer",
	// 		Value: user.Name,
	// 	},
	// }
	// Github API authorization
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)
	if strings.Contains(text, "pr") || strings.Contains(text, "pull request") {
		// Get the list of PRS and their
		message := ""
		if strings.Contains(text, "celts") {
			message += "The current PRs in celts are: \n" + GetPRs("celts", githubClient, ctx)
		} else if strings.Contains(text, "lsf") {
			message += "The current PRs in lsf are: \n" + GetPRs("lsf", githubClient, ctx)
		} else if strings.Contains(text, "bcsr") {
			message += "The current PRs in bcsr are: \n" + GetPRs("bcsr", githubClient, ctx)
		} else {
			message += "The current active PRs for the SSDT are: \n" + GetPRs("all", githubClient, ctx)
		}
		attachment.Text = fmt.Sprintf("Hello %s. %s.", user.RealName, message)

		attachment.Color = "#4af030"
	} else if strings.Contains(text, "weather") {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("Weather is sunny today. %s", user.Name)
		// attachment.Pretext = "How can I be of service"
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("I am good. How are you %s?", user.Name)
		// attachment.Pretext = "How can I be of service"
		attachment.Color = "#4af030"
	}
	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err = slackClient.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
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
