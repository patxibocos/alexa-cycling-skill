package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"os"
)

func HandleRequest(ctx context.Context) (string, error) {
	appConfigUrl := os.Getenv("AWS_APPCONFIG_URL")
	response, err := http.Get(appConfigUrl)
	if err != nil {
		return "", err
	}
	bytes, _ := ioutil.ReadAll(response.Body)
	cyclingData := new(pcsscraper.CyclingData)
	_ = proto.Unmarshal(bytes, cyclingData)
	whatever := cyclingData.Riders[0].LastName
	return fmt.Sprintf("Hi there %s", whatever), nil
}

func main() {
	lambda.Start(HandleRequest)
}
