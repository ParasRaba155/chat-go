#!/bin/bash

read -p "Enter the file name: " file
migrate create -ext sql -dir migrate/migrations -seq $file

