package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Metadata Metadata `json:"metadata"`
	events.APIGatewayProxyResponse
}

type Metadata struct {
	Version         int  `json:"version"`
	BuilderFunction bool `json:"builder_function"`
}

func handler(request events.APIGatewayProxyRequest) (*Response, error) {
	fmt.Printf("%#v\n", request)
	return nil, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
