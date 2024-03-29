FROM golang:1.19.4-alpine
RUN apk update && apk add build-base cmake gcc git
WORKDIR /go/src/github/kardiachain/kai-watcher
ADD . .
WORKDIR /go/src/github/kardiachain/kai-watcher/cmd
RUN go install
WORKDIR /go/bin
