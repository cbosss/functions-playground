package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pkg/errors"

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

var urlRegex = regexp.MustCompile("(.*)/width/(.*)")

func handler(request events.APIGatewayProxyRequest) (*Response, error) {
	matches := urlRegex.FindStringSubmatch(request.Path)

	if len(matches) != 3 {
		return nil, errors.New(fmt.Sprintf("invalid path: %s", request.Path))
	}

	u := url.URL{
		Scheme: "https",
		Path:   request.Path,
		Host:   request.Headers["host"],
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed getting original")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading original body")
	}

	return &Response{
		Metadata: Metadata{
			Version:         1,
			BuilderFunction: true,
		},
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            base64.StdEncoding.EncodeToString(body),
			IsBase64Encoded: true,
		},
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
