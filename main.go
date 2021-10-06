package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

const failHeader = "x-nf-fail"
const statusCodeHeader = "x-nf-status-code"

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.Headers[failHeader] == "fail" {
		return nil, errors.New("fail header detected")
	}

	status := http.StatusOK
	if request.Headers[statusCodeHeader] != "" {
		statusi, err := strconv.ParseInt(request.Headers[statusCodeHeader], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed processing ")
		}
		status = int(statusi)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      status,
		Headers:         map[string]string{"content-type": "application/json"},
		Body:            fmt.Sprintf(`{"timestamp": %s}`, time.Now()),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
