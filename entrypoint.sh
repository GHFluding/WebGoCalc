#!/bin/bash

echo "Waiting for PostgreSQL to be ready..."

# wait postgres database creating
/wait-for-it.sh $DB_HOST:$DB_PORT --timeout=30 --strict -- echo "PostgreSQL is up and ready"
echo "PostgreSQL container is ready."

# start app
echo "Starting the app..."
exec "$@"
