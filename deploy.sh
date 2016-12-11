#!/usr/bin/env bash

source ./common.sh

log INFO "docker compose begin ..."
docker-compose pull
docker-compose -p zqc-server-demo up -d
log INFO "docker compose end"

log INFO "create db indexes begin ..."
docker-compose -p zqc-server-demo exec server zqc db createIndexes
log INFO "create db indexes end"

exit 0
