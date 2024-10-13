#!/bin/bash

set -e
set -o pipefail

COV_FILE_DIR="/tmp"
COV_FILE_NAME="test.cov"
FULL_COV_FILE_NAME="$COV_FILE_DIR/$COV_FILE_NAME"

# Run Go tests with coverage for all packages
# and generate the coverage file
go test -covermode=count \
  -coverprofile="$FULL_COV_FILE_NAME" \
  -coverpkg=./... \
  ./...

# Check if tests pass
if [ $? -ne 0 ]
then
  echo "[ERROR] test execution failed!"
  exit 1
fi

echo "Coverage file generated as $FULL_COV_FILE_NAME"
