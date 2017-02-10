#!/usr/bin/env bash

./deploy-test.sh

log INFO "run unittest begin ..."
docker-compose -p zqc-server-demo-test exec server goconvey -launchBrowser=false
log INFO "run unittest end"
