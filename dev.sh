#!/usr/bin/env sh

# Compile application and run with provided arguments
set -e

PROGRAM_NAME="ovpnller"

go build -o "$PROGRAM_NAME" cmd/ovpnller/*.go

./"$PROGRAM_NAME" "$@"

