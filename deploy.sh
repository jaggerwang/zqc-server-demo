#!/usr/bin/env bash

source ./common.sh

log INFO "deploy development begin ..."
docker-compose pull
docker-compose -p zqc-server-demo up -d
log INFO "deploy development end"

log INFO "create db indexes begin ..."
docker-compose -p zqc-server-demo exec server zqc db createIndexes
log INFO "create db indexes end"

exit 0
