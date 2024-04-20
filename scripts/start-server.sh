#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

go run github.com/joho/godotenv/cmd/godotenv@latest \
    -f $SCRIPT_DIR/../.env \
    go run $SCRIPT_DIR/../src/main.go
