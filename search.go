package main

import (
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-lambda-go/events"
  "github.com/guregu/dynamo"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/endpoints"
  "fmt"
  "os"
  "encoding/json"
)

type AddressData struct {
  Zipcode string `dynamo:"zipcode"`
  Address string
  IsOffce bool
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  db := dynamo.New(session.New(), &aws.Config{
    Region: aws.String(endpoints.ApNortheast1RegionID),
  })
  var tableName = os.Getenv("TABLE")
  table := db.Table(tableName)

  zipcode := request.QueryStringParameters["zipcode"]
  fmt.Println("zipcode: ", zipcode)

  var result AddressData
  err := table.Get("zipcode", zipcode).One(&result)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  res, err := json.Marshal(result)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  return events.APIGatewayProxyResponse{Body: string(res), StatusCode: 200}, nil
}

func main() {
  lambda.Start(Handler)
}
