#!/bin/bash

cd cf-service-connect-repo || exit

go test ./...
