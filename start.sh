#!/bin/bash
export SERVER_ADDRESS=localhost
export SERVER_PORT=8000
export DB_USER=root
export DB_PASSWD=root
export DB_ADDR=localhost
export DB_PORT=3306
export DB_NAME=users
export GOOGLE_AUTH_CLIENT_ID=594665614368-76ld992rnjnlj74o79lalajb3c3u5dff.apps.googleusercontent.com
export GOOGLE_AUTH_CLIENT_SECRET=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
export GOOGLE_AUTH_REDIRECT_URL=http://localhost:8000/callback
go run main.go
