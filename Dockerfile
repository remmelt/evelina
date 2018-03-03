FROM golang:1.9-alpine AS builder
ADD . /go/src/github.com/remmelt/evelina/
WORKDIR /go/src/github.com/remmelt/evelina/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM scratch
WORKDIR /
COPY --from=0 /go/src/github.com/remmelt/evelina/main .
ENTRYPOINT ["/main"]
