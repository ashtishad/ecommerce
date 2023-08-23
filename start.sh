#!/bin/bash
export SERVER_ADDRESS=localhost
export SERVER_PORT=8000
export DB_USER=root
export DB_PASSWD=root
export DB_ADDR=localhost
export DB_PORT=3306
export DB_NAME=users
export GOOGLE_AUTH_REDIRECT_URL=http://localhost:8000/google-auth/callback
go run main.go
