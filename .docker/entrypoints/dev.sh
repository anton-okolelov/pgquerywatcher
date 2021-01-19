#!/usr/bin/env sh

alias target_db="PGPASSWORD=\$TARGET_DB_PASSWORD psql -U \$TARGET_DB_USER -h \$TARGET_DB_HOST -p \$TARGET_DB_PORT \$TARGET_DB_DATABASE"

cd cmd
go run main.go watch \
  --target_db_user $TARGET_DB_USER \
  --target_db_database $TARGET_DB_DATABASE \
  --target_db_host $TARGET_DB_HOST \
  --target_db_port $TARGET_DB_PORT