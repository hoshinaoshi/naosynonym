#!/bin/sh
GOOS=linux go build -o ./bin/main
zip ./bin/main.zip ./bin/main
