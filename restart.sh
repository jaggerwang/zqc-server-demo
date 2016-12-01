#!/usr/bin/env bash

source ./common.sh

./build.sh
if [[ $? != 0 ]]; then
	exit 1
fi

supervisorctl restart all

exit 0
