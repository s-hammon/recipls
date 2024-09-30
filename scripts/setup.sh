#!/bin/bash

if ! go version &> /dev/null; then
    echo "Go is not installed in your environment."
    exit 1
fi

if ! sqlc version &> /dev/null; then
    echo "sqlc is not installed. Installing..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
else
    echo "sqlc is already installed."
fi

if ! goose -version &> /dev/null; then
    echo "goose is not installed. Installing..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
else
    echo "goose is already installed"
fi

ENV_FILE=".env"

if [ -f "$ENV_FILE" ]; then
    echo "$ENV_FILE already exists."
    exit 0
fi

if [ -z "$1" ]; then
    echo "Usage $0 <connection_string>"
    exit 1
fi

{
    echo "DATABASE_URL=\"$1\""
    echo "SCHEMA_DIR=sql/schema"
    echo "DRIVER=postgres"
    echo "JWT_SECRET=\"$(openssl rand -base64 64)\""
} > "$ENV_FILE"

echo "$ENV_FILE created. Please review before running application."