#!/bin/bash
# entrypoint.sh


# Выполняем миграции
echo "Running migrations..."
migrate -path=/app/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

# После выполнения миграций запускаем приложение
echo "Starting the app..."
exec "$@"
