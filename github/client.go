package github

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"

	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func createClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	client.UserAgent = "evelina/0x000"

	return client
}

func sendComment(body string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	comment := github.RepositoryComment{Body: &body}
	s, _, err := client.Repositories.CreateComment(context.Background(), payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, &comment)
	if err != nil {
		return err
	}
	l.WithFields(log.Fields{
		"commentId":  *s.ID,
		"commentURL": *s.URL,
	}).Debug("comment created")
	return nil
}

func setStatusPending(description string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	state := "pending"
	return setStatus(state, description, payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, client, l)
}

func setStatusSuccess(description string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	state := "success"
	return setStatus(state, description, payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, client, l)
}

func setStatusError(description string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	state := "error"
	return setStatus(state, description, payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, client, l)
}

func setStatusFailure(description string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	state := "failure"
	return setStatus(state, description, payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, client, l)
}

func setStatus(state string, description string, payloadRepoOwnerLogin string, payloadRepoName string, prHeadSHA string, client *github.Client, l *log.Entry) error {
	statusContext := "eve"
	status := github.RepoStatus{State: &state, Description: &description, Context: &statusContext}
	s, _, err := client.Repositories.CreateStatus(context.Background(), payloadRepoOwnerLogin, payloadRepoName, prHeadSHA, &status)
	if err != nil {
		return err
	}
	l.WithFields(log.Fields{
		"statusId":  *s.ID,
		"statusURL": *s.URL,
	}).Debug("status created")
	return nil
}

func getPullRequest(payloadRepoOwnerLogin string, payloadRepoName string, payloadIssueNumber int, client *github.Client, l *log.Entry) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Get(context.Background(), payloadRepoOwnerLogin, payloadRepoName, payloadIssueNumber)
	if err != nil {
		l.WithFields(log.Fields{"pullRequestNumber": payloadIssueNumber}).Info("Could not get info for PR")
		return nil, err
	}
	if pr == nil {
		l.WithFields(log.Fields{"pullRequestNumber": payloadIssueNumber}).Info("PR was nil", payloadIssueNumber)
		return nil, fmt.Errorf("PR returned from API was nil")
	}
	return pr, nil
}

// TODO now passing pr object, but we may not need all of it to start the test run.
func callTests(delivery string, pr *github.PullRequest, payloadRepoOwnerLogin string, payloadRepoName string, payloadIssueNumber int, client *github.Client, l *log.Entry) error {
	if pr == nil {
		var err error
		pr, err = getPullRequest(payloadRepoOwnerLogin, payloadRepoName, payloadIssueNumber, client, l)
		if err != nil {
			return err
		}
	}

	if err := setStatusPending(fmt.Sprintf("evelina is running tests — %s", delivery),
		payloadRepoOwnerLogin, payloadRepoName, *pr.Head.SHA, client, l); err != nil {
		return err
	}

	// TODO make %[1] index
	body := fmt.Sprintf("delivery: %s \n"+
		"[kibana](http://kibana/search/for/%s)\n"+
		"[nomad](http://nomad/deployment/%s",
		delivery, delivery, delivery)
	if err := sendComment(body, payloadRepoOwnerLogin, payloadRepoName, *pr.Head.SHA, client, l); err != nil {
		return err
	}

	// TODO call actual test mechanism
	time.Sleep(10 * time.Second)

	if err := setStatusSuccess("tests complete", payloadRepoOwnerLogin, payloadRepoName, *pr.Head.SHA, client, l); err != nil {
		return err
	}
	return nil
}
