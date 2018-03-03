package main

import (
	"bytes"
	"encoding/json"
	"github.com/remmelt/evelina/github"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ghDeliveryHeader = "X-GitHub-Delivery"
const ghEventHeader = "X-GitHub-Event"

func handlePullRequestOpened(delivery string, payload github.PayloadPullRequestOpened) error {
	l(delivery, "Handling PR opened")

	l(delivery, payload.PullRequest.Url)

	return nil
}

func handleIssueCommentCreated(delivery string, payload github.PayloadIssueCommentCreated) error {
	l(delivery, "Handling issue created")

	if !strings.Contains(payload.Comment.Body, "test this") {
		l(delivery, "Not triggering, no match", payload.Comment.Body)
		return nil
	}

	l(delivery, "Trigger a test run for PR", payload.Issue.Number)

	return nil
}

func l(delivery string, msg1 interface{}, msg ...interface{}) {
	log.Println(delivery, msg1, msg)
}

func pr(msg ...interface{}) {
	log.Println(msg)
}

func serve(responseWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != "POST" {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		pr("Received message that was not POST, discarding")
		return
	}

	event := req.Header.Get(ghEventHeader)
	delivery := req.Header.Get(ghDeliveryHeader)
	if event == "" || delivery == "" {
		pr("event or delivery header not found, discarding")
		return
	}

	//pr(ghEventHeader, event)
	//pr(ghDeliveryHeader, delivery)
	//t, _ := httputil.DumpRequest(req, true)
	//pr(string(t))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l("Could not read request body, discarding", err)
		return
	}

	var payload github.GenericPayload
	if err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload); err != nil {
		l("Could not decode payload, discarding", err)
	}

	switch event {
	case "pull_request":
		switch payload.Action {
		case "opened":
			var payload github.PayloadPullRequestOpened
			if err = decode(&payload, body); err == nil {
				handlePullRequestOpened(delivery, payload)
			}
		default:
			l(delivery, "Handling pull_request/"+payload.Action+" not implemented")
		}
	case "issue_comment":
		switch payload.Action {
		case "created":
			var payload github.PayloadIssueCommentCreated
			if err = decode(&payload, body); err == nil {
				handleIssueCommentCreated(delivery, payload)
			}
		default:
			l(delivery, "Handling issue_comment/"+payload.Action+" not implemented")
		}
	default:
		l(delivery, "Handling event "+event+" not implemented")
	}

	responseWriter.Write([]byte(delivery))
}

func decode(payload interface{}, body []byte) error {
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&payload); err != nil {
		l("Could not decode payload, discarding", err)
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/hook/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
