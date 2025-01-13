FROM golang:1.23.3-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk update && apk add --no-cache git wget bash

# Скачиваем migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

WORKDIR /api/backend

ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main ./cmd/apk/main.go

# Скопировать файл в рабочую директорию контейнера
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Сделать файл исполнимым
RUN chmod +x /usr/local/bin/wait-for-it.sh
# Копируем entrypoint-скрипт в контейнер
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Устанавливаем entrypoint
ENTRYPOINT ["/entrypoint.sh", "./main"]
