package main

import (
	"bytes"
	"encoding/json"
	gh "github.com/google/go-github/github"
	"github.com/remmelt/evelina/github"
	l "github.com/remmelt/evelina/util"
	"io/ioutil"
	"log"
	"net/http"
)

const ghDeliveryHeader = "X-GitHub-Delivery"
const ghEventHeader = "X-GitHub-Event"

func serve(responseWriter http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != "POST" {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		l.Pr("Received message that was not POST, discarding")
		return
	}

	event := req.Header.Get(ghEventHeader)
	delivery := req.Header.Get(ghDeliveryHeader)
	if event == "" || delivery == "" {
		l.Pr("event or delivery header not found, discarding")
		return
	}

	//pr(ghEventHeader, event)
	//pr(ghDeliveryHeader, delivery)
	//t, _ := httputil.DumpRequest(req, true)
	//pr(string(t))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l.L(delivery, "Could not read request body, discarding", err)
		return
	}

	switch event {
	case "push":
		var payload gh.PushEvent
		if err = decode(delivery, &payload, body); err == nil {
			github.HandlePush(delivery, payload)
		}
	case "pull_request":
		var payload gh.PullRequestEvent
		if err = decode(delivery, &payload, body); err == nil {
			github.HandlePullRequestOpened(delivery, payload)
		}
	case "issue_comment":
		var payload gh.IssueCommentEvent
		if err = decode(delivery, &payload, body); err == nil {
			github.HandleIssueCommentCreated(delivery, payload)
		}
	default:
		l.L(delivery, "Handling event "+event+" not implemented")
	}

	responseWriter.Write([]byte(delivery))
}

func decode(delivery string, payload interface{}, body []byte) error {
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&payload); err != nil {
		l.L(delivery, "Could not decode payload, discarding", err)
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
