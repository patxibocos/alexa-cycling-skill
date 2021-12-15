package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/alexa"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
)

var s3Client *s3.S3

func getS3Client() *s3.S3 {
	if s3Client == nil {
		sess := session.Must(session.NewSession())
		s3Client = s3.New(sess)
	}
	return s3Client
}

func IntentDispatcher(request alexa.Request, cyclingData *pcsscraper.CyclingData) alexa.Response {
	if request.Body.Intent.Name == "RaceResult" {
		raceNameSlot := request.Body.Intent.Slots["raceName"]
		raceId := raceNameSlot.SlotValue.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
		var race *pcsscraper.Race
		for _, r := range cyclingData.Races {
			if r.Id == raceId {
				race = r
			}
		}
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
		message := ""
		switch ri := raceResult.(type) {
		case *cycling.PastRace:
			message = fmt.Sprintf("Winner was %s %s", ri.GcTop3[0].FirstName, ri.GcTop3[0].LastName)
		case *cycling.FutureRace:
			message = fmt.Sprintf("Race %s happens on %s", race.Name, race.StartDate.AsTime().Format("02-01-2006"))
		case *cycling.RestDayStage:
			message = "RestDayStage"
		case *cycling.SingleDayRaceWithResults:
			message = "SingleDayRaceWithResults"
		case *cycling.SingleDayRaceWithoutResults:
			message = "SingleDayRaceWithoutResults"
		case *cycling.MultiStageRaceWithResults:
			message = "MultiStageRaceWithResults"
		case *cycling.MultiStageRaceWithoutResults:
			message = "MultiStageRaceWithoutResults"
		}
		return alexa.Response{
			Body: alexa.ResBody{
				OutputSpeech: &alexa.OutputSpeech{
					Text: message,
				},
			},
		}
	}
	return alexa.Response{}
}

func Handler(request alexa.Request) (alexa.Response, error) {
	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	s3ObjectKey := os.Getenv("AWS_S3_OBJECT_KEY")
	output, _ := getS3Client().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3ObjectKey),
	})
	cyclingData := new(pcsscraper.CyclingData)
	body, _ := ioutil.ReadAll(output.Body)
	_ = proto.Unmarshal(body, cyclingData)
	return IntentDispatcher(request, cyclingData), nil
}

func main() {
	lambda.Start(Handler)
}
