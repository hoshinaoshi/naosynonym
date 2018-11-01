#!/bin/sh
aws --endpoint-url=http://localhost:4569 dynamodb create-table --table-name test --attribute-definitions AttributeName=tag,AttributeType=S --key-schema AttributeName=tag,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
aws --endpoint-url=http://localhost:4569 dynamodb put-item --table-name test --item '{ "tag": { "S": "testtest" } }'
