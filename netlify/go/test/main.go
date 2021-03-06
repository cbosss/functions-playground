package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

const (
	failHeader       = "x-nf-should-fail"
	notBuilderHeader = "x-nf-not-builder"
	statusCodeHeader = "x-nf-status-code"
	staleAtHeader    = "x-nf-stale-at"
	freshForHeader   = "x-nf-fresh-for"
	headerSizeKB     = "x-nf-header-size-kb"
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
	if request.Headers[failHeader] != "" {
		return nil, errors.New("fail header detected")
	}

	status := http.StatusOK
	if request.Headers[statusCodeHeader] != "" {
		statusi, err := strconv.ParseInt(request.Headers[statusCodeHeader], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed processing")
		}
		status = int(statusi)
	}

	builder := true
	if request.Headers[notBuilderHeader] != "" {
		builder = false
	}

	// if this fails it will just be 0
	size, _ := strconv.ParseInt(request.Headers[headerSizeKB], 10, 64)
	headers := generateHeaders(int(size))
	headers["content-type"] = "application/json"

	if v := request.Headers[freshForHeader]; v != "" {
		if v == "invalid" {
			headers[staleAtHeader] = "invalid"
		}

		if dur, err := time.ParseDuration(v); err == nil {
			timestamp := time.Now().Add(dur)
			headers[staleAtHeader] = strconv.FormatInt(timestamp.Unix(), 10)
		}
	}

	return &Response{
		Metadata: Metadata{
			Version:         1,
			BuilderFunction: builder,
		},
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			StatusCode:      status,
			Headers:         headers,
			Body:            fmt.Sprintf(`{"timestamp": %s}`, time.Now()),
			IsBase64Encoded: false,
		},
	}, nil
}

func generateHeaders(n int) map[string]string {
	headers := make(map[string]string)
	for i := 0; i < n; i++ {
		headers[strconv.FormatInt(int64(i), 10)] = randString(1024)
	}
	return headers
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
