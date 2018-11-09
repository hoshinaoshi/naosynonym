package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  //"github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/guregu/dynamo"
)

type Event struct {
  tag string `json:"name"`
}

type Response struct {
  Result string `json:"tag:"`
}


type Test struct {
  tag string `dynamo:"tag"`
}

func synonym(event Event) (Response, error) {
  //ep := "http://localhost:4569"
  //cred := credentials.NewStaticCredentials("dumy", "dumy", "")
  println("pre session")
  region := "us-west-2"
  conf := &aws.Config{
    //Credentials: cred,
    Region:      &region,
    //Endpoint:    &ep,
  }
  sess, err := session.NewSession(conf)
  if err != nil {
    panic(err)
  }
  println("post session")

  db := dynamo.New(sess)
  table := db.Table("dev-synonyms")

  var result Test
  err = table.Get("tag", "testtest").One(&result)

  println(result.tag)

  return Response{Result: result.tag}, err
}
func main(){
  //synonym(Event{tag: "aa"})
  lambda.Start(synonym)
}
