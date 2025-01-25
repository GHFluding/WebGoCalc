FROM golang:1.23.3-alpine AS builder

# Download package
RUN apk update && apk add --no-cache git wget bash

# Download migrate
RUN wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

WORKDIR /app

ENV GOPROXY=direct
COPY ./migrations /app/migrations

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main ./cmd/apk/main.go
COPY ./docs /docs

COPY ./test/integration_test.go /test/integration_test.go
RUN go test -c -o app.test ./test/integration_test.go

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

#test container
FROM alpine:latest AS test-stage-runner

RUN apk update && apk add --no-cache bash libpq

WORKDIR /test

COPY --from=builder /app/app.test /test/app.test
COPY --from=builder /app/migrations /test/migrations
COPY --from=builder /app/.env /test/.env

# Запуск тестов
CMD ["./app.test", "-test.v"]



RUN touch /test/test_completed





# Этап с минимальным образом для исполнения
FROM alpine:latest AS runner

# Устанавливаем необходимые зависимости для выполнения
RUN apk update && apk add --no-cache bash libpq

# Копируем исполнимая программа и другие необходимые файлы

COPY --from=builder /app/main /main
COPY --from=builder /entrypoint.sh /entrypoint.sh
COPY --from=builder /app/.env /.env
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /docs /docs
# Делаем скрипты исполнимыми
RUN chmod +x  /entrypoint.sh

# Устанавливаем рабочую директорию
WORKDIR /

# Сделаем entrypoint для контейнера
ENTRYPOINT ["/entrypoint.sh", "./main"]