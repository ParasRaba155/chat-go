#!/bin/bash

echo "Installing air"
go install github.com/cosmtrek/air@latest

echo "Installing sqlc"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

echo "Installing migrate"
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
