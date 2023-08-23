#!/bin/bash

TAG=$(git describe --tags)
echo "$TAG" > tag

LAST_TAG=$(git describe --tags --abbrev=0)
git log "$LAST_TAG..HEAD" --oneline --no-decorate > releaselog.md
