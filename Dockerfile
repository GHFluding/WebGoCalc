FROM golang:1.23.3-alpine AS builder

# Download package
RUN apk update && apk add --no-cache git wget bash

# Download migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

WORKDIR /api/backend

ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main ./cmd/apk/main.go

# copy scripts and run
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh







# Этап с минимальным образом для исполнения
FROM alpine:latest

# Устанавливаем необходимые зависимости для выполнения
RUN apk update && apk add --no-cache bash libpq

# Копируем исполнимая программа и другие необходимые файлы
COPY --from=builder /api/backend/main /main
COPY --from=builder /wait-for-it.sh /wait-for-it.sh
COPY --from=builder /entrypoint.sh /entrypoint.sh
COPY --from=builder /api/backend/.env /.env

# Делаем скрипты исполнимыми
RUN chmod +x /wait-for-it.sh /entrypoint.sh

# Устанавливаем рабочую директорию
WORKDIR /

# Сделаем entrypoint для контейнера
ENTRYPOINT ["/entrypoint.sh", "./main"]