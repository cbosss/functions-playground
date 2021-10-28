package main

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const pauseHeader = "x-nf-pause"

type Response struct {
	Metadata Metadata `json:"metadata"`
	events.APIGatewayProxyResponse
}

type Metadata struct {
	Version         int  `json:"version"`
	BuilderFunction bool `json:"builder_function"`
}

type Body struct {
	Start time.Time
	End   time.Time
}

func (b Body) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Start     string
		End       string
		StartUnix int64
		EndUnix   int64
	}{
		b.Start.Format(time.RFC3339),
		b.End.Format(time.RFC3339),
		b.Start.Unix(),
		b.End.Unix(),
	})
}

func handler(request events.APIGatewayProxyRequest) (*Response, error) {
	start := time.Now()

	if pause := request.Headers[pauseHeader]; pause != "" {
		dur, err := time.ParseDuration(pause)
		if err != nil {
			return nil, errors.Wrap(err, "failed parsing duration")
		}
		time.Sleep(dur)
	}

	body, err := json.Marshal(Body{
		Start: start,
		End:   time.Now(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed marshaling body")
	}

	return &Response{
		Metadata: Metadata{
			Version:         1,
			BuilderFunction: true,
		},
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            string(body),
			IsBase64Encoded: false,
		},
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
