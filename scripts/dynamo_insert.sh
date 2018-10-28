#!/bin/sh
aws --endpoint-url=http://localhost:4569 dynamodb create-table --table-name test --attribute-definitions AttributeName=Name,AttributeType=S --key-schema AttributeName=Name,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
aws --endpoint-url=http://localhost:4569 dynamodb put-item --table-name test --item '{ "Name": { "S": "testtest" } }'
