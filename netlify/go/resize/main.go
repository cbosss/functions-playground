package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"golang.org/x/image/draw"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

type Response struct {
	Metadata Metadata `json:"metadata"`
	events.APIGatewayProxyResponse
}

type Metadata struct {
	Version         int  `json:"version"`
	BuilderFunction bool `json:"builder_function"`
}

var urlRegex = regexp.MustCompile("(.*)/ratio/(.*)")

func handler(request events.APIGatewayProxyRequest) (*Response, error) {

	matches := urlRegex.FindStringSubmatch(request.Path)

	if len(matches) != 3 {
		return nil, errors.New(fmt.Sprintf("invalid path: %s", request.Path))
	}

	ratio, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing ratio")
	}

	u := url.URL{
		Scheme: "https",
		Path:   matches[1],
		Host:   request.Headers["host"],
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed getting original")
	}
	defer resp.Body.Close()

	// Decode the image (from PNG to image.Image):
	src, err := png.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode original image")
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/int(ratio), src.Bounds().Max.Y/int(ratio)))

	// resize
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	body := bytes.NewBuffer(nil)
	err = png.Encode(body, dst)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode")
	}

	return &Response{
		Metadata: Metadata{
			Version:         1,
			BuilderFunction: true,
		},
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            base64.StdEncoding.EncodeToString(body.Bytes()),
			IsBase64Encoded: true,
		},
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
