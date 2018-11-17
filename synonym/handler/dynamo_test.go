package handler

import (
  "log"
  "fmt"
  "testing"
  "encoding/json"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/guregu/dynamo"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type TestStruct struct {
    Name string `json:"name"`
}

type Synonym struct {
  Tag        string `dynamo:"tag"`
  Synonyms []string `dynamo:"synonyms,set"`
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

func setCondition() {
  sdkDB := dynamodb.New(session.New(), &aws.Config{
    Region: aws.String("us-west-2"),
    Endpoint: aws.String("http://localhost:4569"),
  })
  packageDB := dynamo.New(session.New(), &aws.Config{
    Region: aws.String("us-west-2"),
    Endpoint: aws.String("http://localhost:4569"),
  })

  if err := packageDB.Table("test-synonyms").DeleteTable().Run(); err != nil {
    fmt.Println()
  }

  input := &dynamodb.CreateTableInput{
    AttributeDefinitions: []*dynamodb.AttributeDefinition{
      {
        AttributeName: aws.String("tag"),
        AttributeType: aws.String("S"),
      },
    },
    KeySchema: []*dynamodb.KeySchemaElement{
      {
        AttributeName: aws.String("tag"),
        KeyType:       aws.String("HASH"),
      },
    },
    ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
      ReadCapacityUnits:  aws.Int64(1),
      WriteCapacityUnits: aws.Int64(1),
    },
    TableName: aws.String("test-synonyms"),
  }
  var err error
  _, err = sdkDB.CreateTable(input)

  if err != nil {
    fmt.Println("Got error calling CreateTable:")
    fmt.Println(err.Error())
  }

  fmt.Println("Created the table test-synonyms in us-west-2")


  table := packageDB.Table("test-synonyms")

  synonym := Synonym{Tag: "testTag", Synonyms: []string{"test_tag", "test-tag"}}
  table.Put(synonym).Run()

  var result Synonym
  table.Get("tag", "testTag").One(&result)
  log.Println(result.Synonyms)
}

func TestHandler(t *testing.T) {
  setCondition()
  queryStringParameters := make(map[string]string)
  queryStringParameters["tag"] = "testTag"
  requestContext := events.APIGatewayProxyRequestContext{Stage: "test"}
  event := events.APIGatewayProxyRequest{QueryStringParameters: queryStringParameters, RequestContext: requestContext}

  response, err := Handler(nil, event)
  if err != nil {
    t.Fatal("fail")
  }
  if response.StatusCode != 200 {
    t.Fatal("fail")
  }
}
