#!/bin/sh
GOOS=linux go build -o ./bin/main
zip ./bin/main.zip ./bin/main
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda create-function --function-name main --runtime go1.x --role r1 --handler main --zip-file fileb://bin/main.zip
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda update-function-code --profile localstack --function-name=main --zip-file fileb://bin/main.zip --publish
