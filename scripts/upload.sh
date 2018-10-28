#!/bin/sh
GOOS=linux go build -o synonym
zip synonym.zip ./synonym
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda create-function --function-name synonym --runtime go1.x --role r1 --handler synonym --zip-file fileb://synonym.zip
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda update-function-code --profile localstack --function-name=synonym --zip-file fileb://synonym.zip --publish
