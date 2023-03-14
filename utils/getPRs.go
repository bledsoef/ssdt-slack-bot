package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func SetupAPI() (context.Context, *github.Client) {
	// set up authentication and api variables
	godotenv.Load(".env")

	githubAccessToken := os.Getenv("GITHUB_AUTH_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)
	return ctx, githubClient
}

func GetAllPRs(githubClient *github.Client, ctx context.Context, owner string) []*github.PullRequest {
	// fetch all of the repositories from the SSDT
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
	return allPRs
}

func ChangesRequired(noReviews bool, pr *github.PullRequest, changeMessage string, lastReview time.Time) string {
	assignee1 := ""
	assignee2 := ""
	if !noReviews {
		assignee1 = "<@" + GetUserID(*pr.Assignees[0].Login) + ">"

		if len(pr.Assignees) > 1 {
			assignee2 = "<@" + GetUserID(*pr.Assignees[1].Login) + ">"

			changeMessage += fmt.Sprintf("%s %s %s: <%s|#%s> %s. Last reviewed on %s. \n", assignee1, assignee2, *pr.GetBase().GetRepo().Name, *pr.HTMLURL, fmt.Sprint(pr.GetNumber()), *pr.Title, lastReview.Format("January 2, 2006"))
		}
	} else {
		assignee1 = "<@" + GetUserID(*pr.Assignees[0].Login) + ">"

		if len(pr.Assignees) > 1 {
			assignee2 = "<@" + GetUserID(*pr.Assignees[1].Login) + ">"

			changeMessage += fmt.Sprintf("%s %s %s: <%s|#%s> %s. Submitted on %s. \n", assignee1, assignee2, *pr.GetBase().GetRepo().Name, *pr.HTMLURL, fmt.Sprint(pr.GetNumber()), *pr.Title, lastReview.Format("January 2, 2006"))
		}
	}
	return changeMessage
}

func ReviewRequired(reviewMessage string, pr *github.PullRequest, lastCommit time.Time) string {
	reviewMessage += fmt.Sprintf("%s: <%s|#%s> %s. Last updated on %s. \n", *pr.GetBase().GetRepo().Name, *pr.HTMLURL, fmt.Sprint(pr.GetNumber()), *pr.Title, lastCommit.Format("January 2, 2006"))
	return reviewMessage
}

func GetOutstandingPRs() string {
	ctx, githubClient := SetupAPI()
	message := ""
	owner := "BCStudentSoftwareDevTeam"
	allPRs := GetAllPRs(githubClient, ctx, owner)

	noReviews := false
	reviewer1, reviewer2 := ReadReviewSchedule()
	reviewMessage := ""

	if reviewer2 != "" {
		reviewMessage = fmt.Sprintf("<@%s> <@%s> *The PRs that require a review are:* \n", reviewer1, reviewer2)

	} else {
		reviewMessage = fmt.Sprintf("<@%s> *The PRs that require a review are:* \n", reviewer1)

	}
	tempReviewMessage := reviewMessage

	changeMessage := "\n*The PRs that require changes are:* \n"
	tempChangeMessage := changeMessage

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

		lastReview := time.Now()

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
			reviewMessage = ReviewRequired(reviewMessage, pr, *lastCommit)
		}

		// if there has recently been a review and still no changes then changes are required
		if !noReviews && ((sinceReview > 48*time.Hour) && (sinceReview < sinceCommit)) {
			changeMessage = ChangesRequired(noReviews, pr, changeMessage, lastReview)
		}
	}
	if reviewMessage == tempReviewMessage {
		reviewMessage = "*There are currently no PRs to review* \n"
	}
	if changeMessage == tempChangeMessage {
		changeMessage = "\n *No PRs currently require changes* \n"
	}
	message += reviewMessage + changeMessage
	return message

}
