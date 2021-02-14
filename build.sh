#!/bin/sh

set -x

GOOS=linux GOARCH=arm GOARM=5 go build

