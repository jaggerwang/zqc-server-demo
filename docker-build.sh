#!/usr/bin/env bash

source ./common.sh

log INFO "docker build begin ..."
docker build -t daocloud.io/jaggerwang/zqc-server-demo .
if [[ $? != 0 ]]; then
	log ERROR "docker build failed"
	exit 1
fi
log INFO "docker build ok"

exit 0
