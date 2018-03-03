package main

import (
	"bytes"
	"encoding/json"
	gh "github.com/google/go-github/github"
	"github.com/remmelt/evelina/github"
	lll "github.com/remmelt/evelina/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const ghDeliveryHeader = "X-GitHub-Delivery"
const ghEventHeader = "X-GitHub-Event"

func serve(responseWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != "POST" {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		lll.Pr("Received message that was not POST, discarding")
		return
	}

	event := req.Header.Get(ghEventHeader)
	delivery := req.Header.Get(ghDeliveryHeader)
	if event == "" || delivery == "" {
		lll.Pr("event or delivery header not found, discarding")
		return
	}

	//pr(ghEventHeader, event)
	//pr(ghDeliveryHeader, delivery)
	//t, _ := httputil.DumpRequest(req, true)
	//pr(string(t))

	l := log.WithFields(log.Fields{"delivery": delivery})

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l.WithError(err).Info("Could not read request body, discarding")
		return
	}

	var e error
	switch event {
	case "push":
		var payload gh.PushEvent
		if e = decode(&payload, body, l); e == nil {
			go github.HandlePush(payload, l)
		}
	case "pull_request":
		var payload gh.PullRequestEvent
		if e = decode(&payload, body, l); e == nil {
			go github.HandlePullRequest(delivery, payload, l)
		}
	case "issue_comment":
		var payload gh.IssueCommentEvent
		if e = decode(&payload, body, l); e == nil {
			go github.HandleIssueComment(delivery, payload, l)
		}
	default:
		//l.Info(delivery, fmt.Sprintf("Handling event '%s' not implemented", event))
	}

	if e != nil {
		l.WithError(e).Info("Error handling request")
	}

	responseWriter.Write([]byte(delivery))
}

func decode(payload interface{}, body []byte, l *log.Entry) error {
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&payload); err != nil {
		l.WithError(err).Info("Could not decode payload, discarding")
		return err
	}
	return nil
}

func main() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatal("No personal access token found. Please provide one as ENV var 'GITHUB_TOKEN'")
	}

	http.HandleFunc("/hook/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
