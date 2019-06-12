#!/usr/bin/env bash
export PORT=3001
export MYSQL_URL="root:root@tcp(127.0.0.1:3306)/dexp?parseTime=true&sql_mode=ansi"
export REDIS_URL=localhost:6379
go run ./server/server.go