VERSION = 0.0.1

# govendor add +external

.PHONY: build
build: evelina

run:
	@gofmt -w -s .
	@go run main.go

evelina:
	@gofmt -w -s .
	@go build -o eve main.go

.PHONY: clean
clean:
	@rm -f eve

.PHONY: package
package: clean
	@docker build -t remmelt/evelina:$(VERSION) .

.PHONY: distribute
distribute: package
	@docker push remmelt/evelina:$(VERSION)

.PHONY: send_pr_opened
send_pr_opened:
	http :8080/hook/ @payloads/pr_opened.json X-Github-Event:pull_request X-Github-Delivery:dddddddd-1eef-11e8-9ea7-0b6cdf45b2db

.PHONY: send_issue_comment_created
send_issue_comment_created:
	http :8080/hook/ @payloads/issue_comment_created.json X-Github-Event:issue_comment X-Github-Delivery:cccccccc-1eef-11e8-9ea7-0b6cdf45b2db

.PHONY: send_issue_comment_created_test_this
send_issue_comment_created_test_this:
	http :8080/hook/ @payloads/issue_comment_created_test_this.json X-Github-Event:issue_comment X-Github-Delivery:bbbbbbbb-1eef-11e8-9ea7-0b6cdf45b2db

.PHONY: send_push_master
send_push_master:
	http :8080/hook/ @payloads/push_master.json X-Github-Event:push X-Github-Delivery:aaaaaaaa-1eef-11e8-9ea7-0b6cdf45b2db
