#!/bin/bash

CURRENT_DIR=$(dirname $(realpath "$0"))

cd "$CURRENT_DIR"
source .env

cd sql/schema
goose postgres $DB_URL $1
