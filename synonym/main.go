package main

import (
  "log"
  "fmt"
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
  Synonyms []*string `json:"synonyms:"`
}

func synonym(_ context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
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
    TableName: aws.String(fmt.Sprintf("%s-synonyms", event.RequestContext.Stage)),
    Key: map[string]*dynamodb.AttributeValue{
      "tag": {
        S: aws.String(request.Tag),
      },
    },
    AttributesToGet: []*string{
      aws.String("synonyms"),
    },
    ConsistentRead: aws.Bool(true),
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
        log.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        log.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        log.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        log.Println(aerr.Error())
      }
    } else {
      log.Println(err.Error())
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

  response := Response{Synonyms: resp.Item["synonyms"].SS}
  responseJson, err := json.Marshal(response)
  if err != nil {
    log.Println("JSON Marshal error:", err)
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
  lambda.Start(synonym)
}
