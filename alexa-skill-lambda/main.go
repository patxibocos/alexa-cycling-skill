package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
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
		whatever = string(bytes)
	}
	return fmt.Sprintf("Hi there %s", whatever), nil
}

func main() {
	lambda.Start(HandleRequest)
}
