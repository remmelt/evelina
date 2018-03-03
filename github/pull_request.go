package github

import (
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

func HandlePullRequestOpened(delivery string, payload github.PullRequestEvent, client *github.Client, l *log.Entry) error {
	l.Info("Handling PR opened")
	l.WithFields(log.Fields{
		"pullRequestNumber": *payload.PullRequest.Number,
		"repoURL":           *payload.Repo.URL,
	}).Info("Trigger a test run for PR")

	if err := callTests(delivery, nil, *payload.Repo.Owner.Login, *payload.Repo.Name, *payload.PullRequest.Number, client, l); err != nil {
		return err
	}

	return nil
}

func HandlePullRequest(delivery string, payload github.PullRequestEvent, l *log.Entry) error {
	client := createClient()

	switch *payload.Action {
	case "opened":
		return HandlePullRequestOpened(delivery, payload, client, l)
	default:
		l.Info("Handling pull_request/", *payload.Action, " not implemented")
	}

	return nil
}
