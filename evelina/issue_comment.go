package evelina

import (
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"strings"
)

func HandleIssueCommentCreated(delivery string, payload github.IssueCommentEvent, client *github.Client, l *log.Entry) error {
	l.Info("Handling issue comment created")

	if !strings.Contains(*payload.Comment.Body, "test this") {
		l.Debug("Not triggering, no match", *payload.Comment.Body)
		return nil
	}

	pr, err := getPullRequest(*payload.Repo.Owner.Login, *payload.Repo.Name, *payload.Issue.Number, client, l)
	if err != nil {
		return err
	}
	if *pr.Merged {
		l.Info(delivery, "This PR is already merged, discarding", *pr.Number)
		return nil
	}

	l.Info("Trigger a test run for PR")
	if err := callTests(delivery, pr, *payload.Repo.Owner.Login, *payload.Repo.Name, *payload.Issue.Number, client, l); err != nil {
		return err
	}

	return nil
}

func HandleIssueComment(delivery string, payload github.IssueCommentEvent, l *log.Entry) error {
	l = l.WithFields(log.Fields{"pullRequestNumber": *payload.Issue.Number, "repoURL": *payload.Repo.URL, "action": *payload.Action})

	client := createClient()

	switch *payload.Action {
	case "created":
		return HandleIssueCommentCreated(delivery, payload, client, l)
	default:
		l.Debug(delivery, "Handling issue_comment/", *payload.Action, " not implemented")
	}

	return nil
}
