#!/bin/bash

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-linux-amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-darwin-amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${OUTPUT_DIR}/credhub-plugin-darwin-arm64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-window-amd64 .
