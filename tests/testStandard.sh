#!/usr/bin/env bash

cd ..

go run main.go \
    -w wordlists/big.txt \
    -u https://www.google.com/FUZZ \
    -fc 403,404 \
    -o tmp/standard.json \
    -X GET
