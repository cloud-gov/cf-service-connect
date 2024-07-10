#!/bin/bash

TAG=$(git describe --tags)
echo "$TAG" > tag

echo "main" > branch-name

# Create release binaries
./bin/create-release-binaries.sh
