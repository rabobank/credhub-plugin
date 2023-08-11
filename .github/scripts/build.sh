#!/bin/bash

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-${VERSION}-linux-amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-${VERSION}-darwin-amd64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${OUTPUT_DIR}/credhub-plugin-${VERSION}-darwin-arm64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_DIR}/credhub-plugin-${VERSION}-window-amd64 .
