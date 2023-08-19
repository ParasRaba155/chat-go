#!/bin/bash

echo "Installing air"
go install github.com/cosmtrek/air@latest

echo "Installing sqlc"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
