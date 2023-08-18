#!/bin/bash
SERVER_ADDRESS=localhost \
SERVER_PORT=8000 \
DB_USER=root \
DB_PASSWD=root \
DB_ADDR=localhost \
DB_PORT=3306 \
DB_NAME=user_auth \
go run main.go
