#!/usr/bin/env bash

cd ./server
GOOS=linux GOARCH=amd64 go build -v -o ../docker/dexp
