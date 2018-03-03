package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type webHookPayload struct {
	Zen    string `json:"zen"`
	HookId uint   `json:"hook_id"`
}

func handlePayload(payload webHookPayload) {
	l(payload.HookId, payload.Zen)

}

func l(hookId uint, msg ...interface{}) {
	log.Printf("%d | %s", hookId, msg)
}

func pr(msg ...interface{}) {
	log.Println(msg)
}

func serve(responseWriter http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		pr("Received message that was not POST, discarding")
		return
	}
	t, _ := httputil.DumpRequest(req, false)
	pr(string(t))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		pr("Could not read request body, discarding", err)
		return
	}
	var payload webHookPayload
	decoder := json.NewDecoder(bytes.NewReader(body))
	err = decoder.Decode(&payload)
	if err != nil {
		pr("Could not decode payload, discarding", err)
	}
	defer req.Body.Close()

	handlePayload(payload)

	responseWriter.Write([]byte(fmt.Sprintf("%d | %s", payload.HookId, payload.Zen)))
}

func main() {
	http.HandleFunc("/hook/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
