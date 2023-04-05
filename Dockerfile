FROM golang:1.17.6-alpine
RUN apk update && apk add build-base cmake gcc git
WORKDIR /go/src/github/kardiachain/kai-watcher
ADD . .
WORKDIR /go/src/github/kardiachain/kai-watcher/cmd/ticker
RUN go install
WORKDIR /go/bin
