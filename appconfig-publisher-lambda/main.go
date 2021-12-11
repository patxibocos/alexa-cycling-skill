package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-sdk-go/service/s3"
)

func HandleRequest(ctx context.Context) (string, error) {
	mySession := session.Must(session.NewSession())
	s3Instance := s3.New(mySession)
	bucket := "alexacycling"
	key := "cycling.data"
	headOutput, err := s3Instance.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "Failed", err
	}
	svc := appconfig.New(mySession)
	applicationId := "lev3aq5"
	configurationProfileId := "izxp027"
	configurationVersion := headOutput.VersionId
	deploymentStrategyId := "AppConfig.AllAtOnce"
	environmentId := "5uih7f1"
	_, err = svc.StartDeployment(&appconfig.StartDeploymentInput{
		ApplicationId:          &applicationId,
		ConfigurationProfileId: &configurationProfileId,
		ConfigurationVersion:   configurationVersion,
		DeploymentStrategyId:   &deploymentStrategyId,
		EnvironmentId:          &environmentId,
	})
	return "Triggered", err
}

func main() {
	lambda.Start(HandleRequest)
}
