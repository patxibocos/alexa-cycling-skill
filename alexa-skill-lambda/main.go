package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
)

func HandleRequest(ctx context.Context) (string, error) {
	response, err := http.Get("http://localhost:2772/applications/alexa_cycling_appconfig/environments/alexa_cycling_appconfig_environment/configurations/alexa_cycling_appconfig_profile")
	var whatever string
	if err != nil {
		whatever = fmt.Sprintf("error: %s", err)
	} else {
		bytes, _ := ioutil.ReadAll(response.Body)
		cyclingData := new(pcsscraper.CyclingData)
		_ = proto.Unmarshal(bytes, cyclingData)
		whatever = cyclingData.Riders[0].LastName
	}
	return fmt.Sprintf("Hi there %s", whatever), nil
}

func main() {
	lambda.Start(HandleRequest)
}
