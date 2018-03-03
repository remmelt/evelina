VERSION = 0.0.1

.PHONY: build
build: evelina

run:
	@gofmt -w -s .
	@go run main.go

evelina:
	@gofmt -w -s .
	@go build -o evelina main.go

.PHONY: clean
clean:
	@rm -f evelina

.PHONY: package
package: clean
	@docker build -t remmelt/evelina:$(VERSION) .

.PHONY: distribute
distribute: package
	@docker push remmelt/evelina:$(VERSION)

.PHONY: send_pr_opened
send_pr_opened:
	http :8080/hook/ @payloads/pr_opened.json X-Github-Event:pull_request X-Github-Delivery:ea0fef50-1eef-11e8-9ea7-0b6cdf45b2db

.PHONY: send_issue_comment_created
send_issue_comment_created:
	http :8080/hook/ @payloads/issue_comment_created.json X-Github-Event:issue_comment X-Github-Delivery:ea0fef50-1eef-11e8-9ea7-0b6cdf45b2db

.PHONY: send_issue_comment_created_test_this
send_issue_comment_created_test_this:
	http :8080/hook/ @payloads/issue_comment_created_test_this.json X-Github-Event:issue_comment X-Github-Delivery:ea0fef50-1eef-11e8-9ea7-0b6cdf45b2db
