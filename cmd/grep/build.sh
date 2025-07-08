#!/bin/bash
go build -o bin/grep-local ./cmd/grep
GOOS=linux GOARCH=amd64 go build -o bin/grep ./cmd/grep