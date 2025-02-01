FROM golang:1.23.3-alpine AS builder

# Download package
RUN apk update && apk add --no-cache git wget bash

# Download migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

WORKDIR /app

ENV GOPROXY=direct
COPY ../../build/migrations /app/migrations/

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../. .


RUN go build -o main ./cmd/apk/main.go
COPY ../../docs/ /docs/


# Этап с минимальным образом для исполнения
FROM scratch AS runner

# Копируем исполнимая программа и другие необходимые файлы

COPY --from=builder /app/main /main
COPY --from=builder /app/.env /.env
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /docs /docs


# Устанавливаем рабочую директорию
WORKDIR /

# Сделаем entrypoint для контейнера
ENTRYPOINT ["./main"]