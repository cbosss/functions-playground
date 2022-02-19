package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Path                  string              `json:"path"`
	SingleValueQueryParam map[string]string   `json:"single_value_query_params"`
	MultiValueQueryParam  map[string][]string `json:"multi_value_query_params"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	payload, err := json.Marshal(Request{
		Path:                  request.Path,
		SingleValueQueryParam: request.QueryStringParameters,
		MultiValueQueryParam:  request.MultiValueQueryStringParameters,
	})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(payload),
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
