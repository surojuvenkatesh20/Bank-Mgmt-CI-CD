#!/bin/sh

set -e  #script will exit immediately if any command returns a non-zero status

echo "run db migrations"
source /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"  #instead of creating a new process for running CMD command from Dockerfile, use the same process PID and then run the CMD commands.