package github

import (
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

func HandlePush(payload github.PushEvent, l *log.Entry) error {
	l.Info("Handling push")

	if *payload.Ref != "refs/heads/master" {
		l.Info("Not triggering, not pushed to master but to", *payload.Ref)
		return nil
	}
	l.WithFields(log.Fields{
		"ref":  *payload.Ref,
		"repo": *payload.Repo.URL,
	}).Info("Trigger a package run, creating a new release")

	return nil
}
