package evelina

import (
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

func HandlePush(payload github.PushEvent, l *log.Entry) error {
	l = l.WithFields(log.Fields{
		"ref":     *payload.Ref,
		"repoURL": *payload.Repo.URL,
		"head":    *payload.Head,
	})
	l.Info("Handling push")

	if *payload.Ref != "refs/heads/master" {
		l.Info("Not triggering, not pushed to master")
		return nil
	}

	return packageApp(payload, l)
}

func packageApp(payload github.PushEvent, l *log.Entry) error {
	l.Info("Trigger a package run, creating a new release")

	if *payload.Repo.URL == "evelina_test"{
		nomad.
	}

	return nil
}
