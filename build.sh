#!/usr/bin/env bash

source ./common.sh

log INFO "build begin ..."
go get -v ./... && go install -v . && mv $GOPATH/bin/zqcserverdemo $GOPATH/bin/zqc
if [[ $? != 0 ]]; then
	log ERROR "build failed"
	exit 1
fi
log INFO "build ok"

exit 0
