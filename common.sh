#!/usr/bin/env bash

log() {
	level=$1
	message=$2
	echo "`hostname` `date +'%Y-%m-%d %H:%M:%S'` $1 $2"
}
