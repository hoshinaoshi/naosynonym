package main

import (
  "fmt"
  "log"
  "context"
  "encoding/json"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type Request struct {
  Tag string `json:"tag"`
}

type Response struct {
  Synonyms string `json:"synonyms:"`
}

func synonym(c context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
  //log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
  log.Printf("Processing Lambda request %s\n", request.QueryStringParameters["tag"])
  log.Printf("Processing Lambda request %s\n", request.QueryStringParameters)
  log.Printf("Processing Lambda request %s\n", request)
  //tag := request.QueryStringParameters["tag"]
  //fmt.Println(tag)dd

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
    response := Response{Synonyms: *resp.Item["tag"].S}
    responseJson, err := json.Marshal(response)
    if err != nil {
      fmt.Println("JSON Marshal error:", err)
    }

    return events.APIGatewayProxyResponse{
      Headers: map[string]string{
        "Content-Type": "application/json",
      },
      StatusCode: 200,
      Body: string(responseJson),
      IsBase64Encoded: false,
    }, err
}
func main(){
  //synonym(Event{name: "aa"})
  lambda.Start(synonym)
}
