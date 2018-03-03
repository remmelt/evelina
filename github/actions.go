package github

import (
	"github.com/google/go-github/github"
	l "github.com/remmelt/evelina/util"
	"strings"
)

func HandlePush(delivery string, payload github.PushEvent) error {
	l.L(delivery, "Handling push")

	if *payload.Ref != "refs/heads/master" {
		l.L(delivery, "Not triggering, not pushed to master but to", *payload.Ref)
		return nil
	}
	l.L(delivery, "Trigger a package run, creating a new release for", *payload.Ref, "on repository", *payload.Repo.URL)

	return nil
}

func HandlePullRequestOpened(delivery string, payload github.PullRequestEvent) error {
	switch *payload.Action {
	case "opened":
		l.L(delivery, "Handling PR opened")
		l.L(delivery, "Trigger a test run for PR", *payload.PullRequest.Number, "on repository", *payload.Repo.URL)
	default:
		l.L(delivery, "Handling pull_request/", *payload.Action, " not implemented")
	}

	return nil
}

func HandleIssueCommentCreated(delivery string, payload github.IssueCommentEvent) error {
	switch *payload.Action {
	case "created":
		l.L(delivery, "Handling issue created")

		if !strings.Contains(*payload.Comment.Body, "test this") {
			l.L(delivery, "Not triggering, no match", *payload.Comment.Body)
			return nil
		}

		l.L(delivery, "Trigger a test run for PR", *payload.Issue.Number, "on repository", *payload.Repo.HTMLURL)
	default:
		l.L(delivery, "Handling issue_comment/", *payload.Action, " not implemented")
	}
	return nil
}
