package main

import (
  "fmt"
  "log"
  "context"
  "encoding/json"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws/awserr"
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

func synonym(c context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
  //log.Printf("Processing Lambda event %s\n", event.RequestContext.RequestID)
  log.Printf("Processing Lambda event %s\n", event.QueryStringParameters["tag"])
  log.Printf("Processing Lambda event %s\n", event.QueryStringParameters)
  log.Printf("Processing Lambda event %s\n", event)

  request := Request{Tag: event.QueryStringParameters["tag"]}
  if request.Tag == "" {
    return events.APIGatewayProxyResponse{
      Headers: map[string]string{
        "Content-Type": "application/json",
      },
      StatusCode: 400,
      Body: "",
      IsBase64Encoded: false,
    }, nil
  }

  ddb := dynamodb.New(session.New())

    params := &dynamodb.GetItemInput{
        TableName: aws.String("dev-synonyms"), // テーブル名

        Key: map[string]*dynamodb.AttributeValue{
            "tag": {             // キー名
                S: aws.String(request.Tag),   // 持ってくるキーの値
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

    if len(resp.Item) == 0 {
      return events.APIGatewayProxyResponse{
        Headers: map[string]string{
          "Content-Type": "application/json",
        },
        StatusCode: 404,
        Body: "",
        IsBase64Encoded: false,
      }, nil
    }
    if err != nil {
      if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case dynamodb.ErrCodeProvisionedThroughputExceededException:
          fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
        case dynamodb.ErrCodeResourceNotFoundException:
          fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
        case dynamodb.ErrCodeInternalServerError:
          fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
        default:
          fmt.Println(aerr.Error())
        }
      } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
      }
      return events.APIGatewayProxyResponse{
        Headers: map[string]string{
          "Content-Type": "application/json",
        },
        StatusCode: 500,
        Body: "error",
        IsBase64Encoded: false,
      }, nil
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
