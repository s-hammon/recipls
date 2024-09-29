#!/bin/bash

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
    echo "CONN_STRING=\"$1\""
    echo "SCHEMA_DIR=sql/schema"
    echo "DRIVER=postgres"
    echo "JWT_SECRET=\"$(openssl rand -base64 64)\""
} > "$ENV_FILE"

echo "$ENV_FILE created. Please review before running application."