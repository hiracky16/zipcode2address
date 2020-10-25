package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"fmt"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	zipcode := request.QueryStringParameters["zipcode"]
  fmt.Println("zipcode: ", zipcode)
  return events.APIGatewayProxyResponse{Body: zipcode, StatusCode: 200}, nil
}

func main() {
  lambda.Start(Handler)
}
