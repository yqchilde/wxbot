#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wxbot
mv ./wxbot ./docker/wxbot && cd ./docker || exit 0
docker buildx build --platform=linux/amd64 -t yqchilde/wxbot .
docker save -o wxbot.tar yqchilde/wxbot
rm -f wxbot
