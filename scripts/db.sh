#!/bin/bash
if [ -f .env ]; then
    source .env
else
    echo ".env file not found. please create or run setup.sh script"
fi

goose -dir ${SCHEMA_DIR} ${DRIVER} "${CONN_STRING}" up 