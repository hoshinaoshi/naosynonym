package handler

import (
  "testing"
  "encoding/json"
  "github.com/aws/aws-lambda-go/events"
)

type TestStruct struct {
    Name string `json:"name"`
}

func TestResponseAPIGatewayProxyResponseSuccess(t *testing.T) {
  testStruct := TestStruct{Name: "test"}
  jsonBytes, _ := json.Marshal(testStruct)

  result, err := ResponseAPIGatewayProxyResponse(jsonBytes, 200, nil)
  if err != nil {
    t.Fatal("response not exist")
  }
  if result.Headers["Content-Type"] != "application/json" {
    t.Fatal("fail content type")
  }
  if result.StatusCode != 200 {
    t.Fatal("fail status code")
  }
  if result.IsBase64Encoded {
    t.Fatal("fail IsBase64Encoded")
  }
}

func TestSynonym(t *testing.T) {
  queryStringParameters := make(map[string]string)
  queryStringParameters["tag"] = "test"
  event := events.APIGatewayProxyRequest{QueryStringParameters: queryStringParameters}

  response, err := Synonym(nil, event)
  if err != nil {
    t.Fatal("fail")
  }
  if response.StatusCode != 200 {
    t.Fatal("fail")
  }
}
