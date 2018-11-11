package main

import (
  "context"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "fmt"
)

type Request struct {
  Tag string `json:"tag"`
}

type Response struct {
  Synonyms string `json:"synonyms:"`
}

func synonym(ctx context.Context, req Request) (Response, error){

  ddb := dynamodb.New(session.New())

    params := &dynamodb.GetItemInput{
        TableName: aws.String("dev-synonyms"), // テーブル名

        Key: map[string]*dynamodb.AttributeValue{
            "tag": {             // キー名
                S: aws.String("testtest"),   // 持ってくるキーの値
            },
        },
        AttributesToGet: []*string{
            aws.String("tag"),     // 欲しいデータの名前
        },
        ConsistentRead: aws.Bool(true),     // 常に最新を取得するかどうか

        //返ってくるデータの種類
        ReturnConsumedCapacity: aws.String("NONE"),
    }

    resp, err := ddb.GetItem(params)

    if err != nil {
        fmt.Println(err.Error())
    }

    //resp.Item[項目名].型 でデータへのポインタを取得
    fmt.Println(*resp.Item["tag"].S)

    return Response{Synonyms: *resp.Item["tag"].S}, err
}
func main(){
  //synonym(Event{name: "aa"})
  lambda.Start(synonym)
}
