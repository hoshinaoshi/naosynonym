package main

import (
  "encoding/csv"
  "fmt"
  "log"
  "os"
  "strings"
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/guregu/dynamo"
)

type Synonym struct {
  Tag        string `dynamo:"tag"`
  Synonyms []string `dynamo:"synonyms,set"`
}

func main() {

  var fp *os.File
  if len(os.Args) < 2 {
    fp = os.Stdin
  } else {
    var err error
    fp, err = os.Open(os.Args[1])
    if err != nil {
      panic(err)
    }
    defer fp.Close()
  }

  reader := csv.NewReader(fp)
  reader.Comma = '\t'

  // コメント設定(なんとコメント文字を指定できる!)
  reader.Comment = '#'

  // 全部読みだす
  records, err := reader.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  db := dynamo.New(session.New(), &aws.Config{
    Region: aws.String("us-west-2"),
  })
  table := db.Table("dev-synonyms")
log.Println(table)
  // 各行でループ
  for _, row := range records {
    fmt.Printf("tag: %s, synonyms: %s", row[0], strings.Split(row[5], "/"))

    synonym := Synonym{Tag: , Synonyms: []string{"b","c"}}
    log.Println(synonym.Synonyms)
    go table.Put(synonym).Run()
  }
  time.Sleep(3 * time.Second)
}
