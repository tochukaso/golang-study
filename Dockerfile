FROM golang:1.16.3 AS base

ENV GOPATH=/go
ENV GO111MODULE=on

WORKDIR /go/src/app

ADD . /go/src/app
