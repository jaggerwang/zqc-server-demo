#!/usr/bin/env bash

source ./common.sh

log INFO "docker compose begin ..."
docker-compose pull
docker-compose -p zqc-server-demo -f docker-compose.yml up -d
log INFO "docker compose end"

log INFO "create db indexes begin ..."
docker-compose -p zqc-server-demo -f docker-compose.yml exec server zqc db createIndexes
log INFO "create db indexes end"

exit 0
