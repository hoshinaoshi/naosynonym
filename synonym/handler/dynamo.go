package handler

import (
  "log"
  "fmt"
  "context"
  "encoding/json"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type Request struct {
  Tag string `json:"tag"`
}

type Response struct {
  Synonyms []*string `json:"synonyms:"`
}

func ResponseAPIGatewayProxyResponse(body []byte, statusCode int, err error) (events.APIGatewayProxyResponse, error) {
  return events.APIGatewayProxyResponse{
    Headers: map[string]string{
      "Content-Type": "application/json",
    },
    StatusCode: statusCode,
    Body: string(body),
    IsBase64Encoded: false,
  }, err
}

func Handler(_ context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
  eventJsonBytes, _ := json.Marshal(event)
  log.Printf("Processing Lambda event %s\n", eventJsonBytes)

  request := Request{Tag: event.QueryStringParameters["tag"]}
  if request.Tag == "" {
    return ResponseAPIGatewayProxyResponse([]byte{}, 400, nil)
  }

  endPoint := ""
  if event.RequestContext.Stage == "test" {
    endPoint = "http://localhost:4569"
  }
  ddb := dynamodb.New(session.New(), &aws.Config{
    Region: aws.String("us-west-2"),
    Endpoint: aws.String(endPoint),
  })

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
    return ResponseAPIGatewayProxyResponse([]byte{}, 404, nil)
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
    return ResponseAPIGatewayProxyResponse([]byte{}, 500, err)
  }

  response := Response{Synonyms: resp.Item["synonyms"].SS}
  responseJson, err := json.Marshal(response)
  if err != nil {
    log.Println("JSON Marshal error:", err)
  }

  return ResponseAPIGatewayProxyResponse(responseJson, 200, nil)
}
