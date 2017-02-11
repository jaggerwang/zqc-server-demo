#!/usr/bin/env bash

source ./common.sh

log INFO "deploy test begin ..."
docker-compose pull
docker-compose -p zqc-server-demo-test -f docker-compose.yml -f docker-compose.test.yml up -d
log INFO "deploy test end"

log INFO "empty db begin ..."
docker-compose -p zqc-server-demo-test exec server zqc db empty
log INFO "empty db end"

log INFO "create db indexes begin ..."
docker-compose -p zqc-server-demo-test exec server zqc db createIndexes
log INFO "create db indexes end"

exit 0
