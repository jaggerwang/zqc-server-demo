#!/usr/bin/env bash

source ./common.sh

log INFO "deploy production begin ..."
docker-compose pull
docker-compose -p zqc-server-demo -f docker-compose.yml -f docker-compose.prod.yml up -d
log INFO "deploy production end"

log INFO "create db indexes begin ..."
docker-compose -p zqc-server-demo exec server zqc db createIndexes
log INFO "create db indexes end"

exit 0
