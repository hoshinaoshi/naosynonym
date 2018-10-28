package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/guregu/dynamo"
)

type Event struct {
  Name string `json:"name"`
}

type Response struct {
  Result string `json:"Name:"`
}


type Test struct {
  Name string `dynamo:"Name"`
}

func synonym(event Event) (Response, error) {
  ep := "http://localhost:4569"
  cred := credentials.NewStaticCredentials("dumy", "dumy", "")
  region := "us-east-1"
  conf := &aws.Config{
    Credentials: cred,
    Region:      &region,
    Endpoint:    &ep,
  }
  sess, err := session.NewSession(conf)
  if err != nil {
    panic(err)
  }

  db := dynamo.New(sess)
  table := db.Table("test")

  var result Test
  err = table.Get("Name", "testtest").One(&result)

  println(result.Name)

  return Response{Result: result.Name}, err
}
func main(){
  //synonym(Event{Name: "aa"})
  lambda.Start(synonym)
}
