#!/bin/bash

CURRENT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$CURRENT_DIR"
. ./initializeEnvironment.sh

goose postgres $CONSTRING $1
