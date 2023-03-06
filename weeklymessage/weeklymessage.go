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
	message := "*Good morning everyone!*"
	var needsReview []*github.PullRequest
	var needsChanges []*github.PullRequest

	// set up authentication and api variables
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)

	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	// fetch all of the repositories from the SSDT
	owner := "BCStudentSoftwareDevTeam"
	repos, _, err := githubClient.Repositories.List(ctx, owner, nil)
	if err != nil {
		fmt.Printf("Failed to list repositories: %v\n", err)
	}

	// compile all of the repos' PRs into one list
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
	noReviews := false
	lastReview := time.Now()
	lastCommit := time.Now()

	// iterate through the prs and identify if they are outstanding.
	for _, pr := range allPRs {
		// Get the list of reviews for the Pull Request
		reviews, _, err := githubClient.PullRequests.ListReviews(ctx, owner, *pr.GetBase().GetRepo().Name, pr.GetNumber(), nil)
		if err != nil {
			fmt.Printf("Failed to list reviews for Pull Request #%d: %v\n", pr.GetNumber(), err)
			continue
		}
		// initialize boolean values to determine if a pr is outstanding or inactice

		noReviews = false

		lastReview = time.Now()

		commits, _, err := githubClient.PullRequests.ListCommits(ctx, owner, *pr.GetBase().GetRepo().Name, pr.GetNumber(), nil)
		if err != nil {
			fmt.Printf("Failed to list commits for Pull Request #%d: %v\n", pr.GetNumber(), err)
			continue
		}
		lastCommit := commits[len(commits)-1].Commit.Committer.Date

		// if a review isn't currently in then set the last review time to be the same as the last
		if len(reviews) > 0 {
			lastReview = reviews[len(reviews)-1].GetSubmittedAt()
		} else {
			lastReview = *lastCommit
			noReviews = true
		}
		sinceReview := time.Since(lastReview)
		sinceCommit := time.Since(*lastCommit)

		// check to see if the PR is inactive
		if (noReviews || (sinceReview > 730*time.Hour)) || (sinceCommit > 730*time.Hour) {
			continue
		}

		// if there have recently been changes and still no review then a review is required
		if (sinceCommit > 48*time.Hour) && ((sinceReview > sinceCommit) || noReviews) {
			needsReview = append(needsReview, pr)

		}

		// if there has recently been a review and still no changes then changes are required
		if !noReviews && ((sinceReview > 48*time.Hour) && (sinceReview < sinceCommit)) {
			needsChanges = append(needsChanges, pr)

		}

		// if needsReview {
		// 	fmt.Printf("Pull Request #%d needs review\n", pr.GetNumber())
		// } else if needsChanges {
		// 	fmt.Printf("Pull Request #%d needs changes\n", pr.GetNumber())
		// }
	}
	// loop throught the prs that require a review and add them to the message
	message += "*The PRs that require a review are:* \n"
	for i, review := range needsReview {
		message += fmt.Sprintf("%s. %s: #%s %s. Last updated on %s. \n", fmt.Sprint(i+1), *review.GetBase().GetRepo().Name, fmt.Sprint(review.GetNumber()), *review.Title, lastCommit.Format("January 2, 2006"))
	}

	// loop throught the prs that require changes and add them to the message
	message += "*The PRs that require changes are:* \n"
	for i, changes := range needsChanges {
		if !noReviews {
			message += fmt.Sprintf("%s. %s: #%s %s. Last reviewed on %s. \n", fmt.Sprint(i+1), *changes.GetBase().GetRepo().Name, fmt.Sprint(changes.GetNumber()), *changes.Title, lastReview.Format("January 2, 2006"))
		} else {
			message += fmt.Sprintf("%s. %s: #%s %s. Submitted on %s. \n", fmt.Sprint(i+1), *changes.GetBase().GetRepo().Name, fmt.Sprint(changes.GetNumber()), *changes.Title, lastCommit.Format("January 2, 2006"))
		}
	}
	return message

}
