#!/bin/bash

PLUGIN_NAME="cf-service-connect"

GOOS=darwin GOARCH=amd64 go build -o ${PLUGIN_NAME}_darwin_amd64
GOOS=linux GOARCH=amd64 go build -o ${PLUGIN_NAME}_linux_amd64
GOOS=linux GOARCH=386 go build -o ${PLUGIN_NAME}_linux_386
GOOS=windows GOARCH=amd64 go build -o ${PLUGIN_NAME}_windows_amd64
GOOS=windows GOARCH=386 go build -o ${PLUGIN_NAME}_windows_386
