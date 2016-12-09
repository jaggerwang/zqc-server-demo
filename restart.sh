#!/usr/bin/env bash

source ./common.sh

log INFO "build begin ..."
go get -d -v ./... && go install -v .
if [[ $? != 0 ]]; then
	log ERROR "build failed"
	exit 1
fi
log INFO "build ok"

supervisorctl restart all

exit 0
