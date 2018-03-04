package evelina

import (
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

func callTestRun(delivery string, payload github.PullRequestEvent, client *github.Client, l *log.Entry) error {
	l.Info("Trigger a test run for PR")

	if err := callTests(delivery, nil, *payload.Repo.Owner.Login, *payload.Repo.Name, *payload.PullRequest.Number, client, l); err != nil {
		return err
	}
	return nil
}

func callCancelTests(delivery string, payload github.PullRequestEvent, client *github.Client, l *log.Entry) error {
	l.Info("Canceling any previous runs")
	return nil
}

func HandlePullRequestOpened(delivery string, payload github.PullRequestEvent, client *github.Client, l *log.Entry) error {
	l.Info("Handling PR opened")
	return callTestRun(delivery, payload, client, l)
}

func HandlePullRequestSynchronize(delivery string, payload github.PullRequestEvent, client *github.Client, l *log.Entry) error {
	l.Info("Handling PR synchronize")

	callCancelTests(delivery, payload, client, l)

	return callTestRun(delivery, payload, client, l)
}

func HandlePullRequest(delivery string, payload github.PullRequestEvent, l *log.Entry) error {
	l = l.WithFields(log.Fields{"pullRequestNumber": *payload.PullRequest.Number, "repoURL": *payload.Repo.URL, "action": *payload.Action})

	client := createClient()

	switch *payload.Action {
	case "opened":
		return HandlePullRequestOpened(delivery, payload, client, l)
	case "synchronize":
		return HandlePullRequestSynchronize(delivery, payload, client, l)
	default:
		l.Info("Handling pull_request/", *payload.Action, " not implemented")
	}

	return nil
}
