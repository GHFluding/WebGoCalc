
FROM golang:1.23.3-alpine AS builder


RUN apk add --no-cache git

WORKDIR /

ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CONFIG_PATH=./config/local.yaml


RUN go build -o main ./cmd/apk/main.go

ENTRYPOINT [ "./main" ]
