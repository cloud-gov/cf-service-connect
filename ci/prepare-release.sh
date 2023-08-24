#!/bin/bash

TAG=$(git describe --tags)
echo "$TAG" > tag

# Create release binaries
./bin/create-release-binaries.sh
