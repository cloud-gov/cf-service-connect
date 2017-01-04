#!/bin/bash

set -e
set -x

NAME=ServiceConnect
SUBCOMMAND=connect-to-service

# http://stackoverflow.com/a/1371283/358804
BIN=${PWD##*/}

go build

# reinstall
if cf plugins | grep -q "$NAME"; then
  cf uninstall-plugin "$NAME"
fi
cf install-plugin -f "$BIN"

DEBUG=1 cf "$SUBCOMMAND" "$@"
