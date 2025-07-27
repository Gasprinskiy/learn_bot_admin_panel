#!/bin/bash

if [[ "$1" == "-prod" ]]; then
  DETACH_FLAG="-d"
else
  DETACH_FLAG=""
fi


./get_pg_env.sh

echo "Run lear_bot_admin_planel"

# Запуск
docker compose down
docker compose build
docker compose up --force-recreate $DETACH_FLAG