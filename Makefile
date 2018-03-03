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
