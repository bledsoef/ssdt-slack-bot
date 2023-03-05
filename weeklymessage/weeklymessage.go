package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"golang.org/x/oauth2"
)

func main() {
	// Set the message text
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")

	channelID := "C04LA97FWKH"

	api := slack.New(token)

	message := getOutstandingPRs()
	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
		return
	}
	// Create a new cron job that runs once a week on Monday at 9am
	// 	cron := cron.New(cron.WithSeconds())
	// 	cron.AddFunc("@daily", func() {
	// 		// Create a time object for the scheduled time
	// 		// Schedule the message to be sent
	// 		_, _, err := api.PostMessage(
	// 			channelID,
	// 			slack.MsgOptionText(message, false),
	// 		)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
	// 			return
	// 		}
	// 	})

	// 	cron.AddFunc("0 0 9 * * 5", func() {
	// 		// Create a time object for the scheduled time
	// 		// Schedule the message to be sent
	// 		_, _, err := api.PostMessage(
	// 			channelID,
	// 			slack.MsgOptionText(message, false),
	// 		)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
	// 			return
	// 		}
	// 	})

	// 	// Start the cron job
	// 	cron.Start()

	// // Wait forever
	// select {}
}

func getOutstandingPRs() string {
	message := ""
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)
	owner := "BCStudentSoftwareDevTeam"
	repos, _, err := githubClient.Repositories.List(ctx, owner, nil)

	if err != nil {
		fmt.Printf("Failed to list repositories: %v\n", err)
	}

	var allPRs []*github.PullRequest
	for _, repo := range repos {
		prs, _, err := githubClient.PullRequests.List(ctx, owner, repo.GetName(), &github.PullRequestListOptions{
			State: "open",
		})
		if err != nil {
			fmt.Printf("Failed to list pull requests for %s: %v\n", repo.GetName(), err)
			continue
		}

		allPRs = append(allPRs, prs...)
	}

	for _, pr := range allPRs {
		// Get the list of reviews for the Pull Request
		reviews, _, err := githubClient.PullRequests.ListReviews(ctx, owner, *pr.GetBase().GetRepo().Name, pr.GetNumber(), nil)
		if err != nil {
			fmt.Printf("Failed to list reviews for Pull Request #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		// Check if any of the reviews have been submitted within the last 48 hours
		needsReview := true
		needsChanges := true
		lastReview := ""
		sinceReview := ""

		if len(reviews) > 0 {
			lastReview := reviews[len(reviews)-1].GetSubmittedAt()

			sinceReview := time.Since(lastReview)
			if (sinceReview < 48*time.Hour) || (730*time.Hour < sinceReview) {
				needsReview = false
			}
			fmt.Println(sinceReview)
			fmt.Println(730 * time.Hour)
			fmt.Println(pr.GetNumber())
		}

		commits, _, err := githubClient.PullRequests.ListCommits(ctx, owner, *pr.GetBase().GetRepo().Name, pr.GetNumber(), nil)
		if err != nil {
			fmt.Printf("Failed to list commits for Pull Request #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		lastCommit := ""
		sinceCommit := ""

		if len(commits) > 0 {
			lastCommit := commits[len(commits)-1].Commit.Committer.Date
			sinceCommit := time.Since(*lastCommit)
			if (sinceCommit < 48*time.Hour) || (730*time.Hour < sinceCommit) {
				needsChanges = false
			}
			fmt.Println(sinceReview)
			fmt.Println(730 * time.Hour)
			fmt.Println(pr.GetNumber())
		}

		if needsChanges && needsReview {
			if (lastReview != "") && (lastCommit != "") {
				if sinceReview <= sinceCommit {
					needsReview = true
				} else if sinceReview > sinceCommit {
					needsChanges = true
				}

			} else if lastReview != "" {
				needsReview = true
			}
		}

		// Mark the Pull Request as needing review if no review has been submitted within the last 24 hours
		// if needsReview {
		// 	fmt.Printf("Pull Request #%d needs review\n", pr.GetNumber())
		// } else if needsChanges {
		// 	fmt.Printf("Pull Request #%d needs changes\n", pr.GetNumber())
		// }
	}
	return message

}
