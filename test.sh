#!/usr/bin/env bash

source ./common.sh

./deploy-test.sh

log INFO "run unittest begin ..."
docker-compose -p zqc-server-demo-test exec server goconvey -host 0.0.0.0 -launchBrowser=false
log INFO "run unittest end"
