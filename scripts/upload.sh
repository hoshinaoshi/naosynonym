#!/bin/sh
GOOS=linux go build -o ./synonym/synonym
zip ./synonym/synonym.zip ./synonym/synonym
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda create-function --function-name ./synonym/synonym --runtime go1.x --role r1 --handler ./synonym/synonym --zip-file fileb://synonym/synonym.zip
aws --endpoint-url=http://localhost:4574 --region us-east-1 --profile localstack lambda update-function-code --profile localstack --function-name=synonym --zip-file fileb://synonym/synonym.zip --publish
