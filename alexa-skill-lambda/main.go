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
	"strings"
	"time"
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
		raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
		var race *pcsscraper.Race
		for _, r := range cyclingData.Races {
			if r.Id == raceId {
				race = r
			}
		}
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
		message := ""
		switch ri := raceResult.(type) {
		case *cycling.PastRace: // We will have a mapping from race id to speakable race name (only in Spanish)
			message = fmt.Sprintf(
				"%s terminó el %s. El ganador fue %s, el segundo %s y tercero %s",
				pcsscraper.RaceName[race.Id],
				formattedDate(race.EndDate.AsTime()),
				cycling.RiderFullName(ri.GcTop3[0]),
				cycling.RiderFullName(ri.GcTop3[1]),
				cycling.RiderFullName(ri.GcTop3[2]),
			)
		case *cycling.FutureRace:
			message = fmt.Sprintf("%s no empieza hasta el %s", pcsscraper.RaceName[race.Id], formattedDate(race.StartDate.AsTime()))
		case *cycling.RestDayStage:
			message = fmt.Sprintf("Hoy ha habido día de descanso en %s", pcsscraper.RaceName[race.Id])
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
			Version: "1.0",
			Body: alexa.ResBody{
				OutputSpeech: &alexa.OutputSpeech{
					Type: "PlainText",
					Text: message,
				},
				ShouldEndSession: true,
			},
		}
	}
	return alexa.Response{
		Version: "1.0",
		Body: alexa.ResBody{
			OutputSpeech: &alexa.OutputSpeech{
				Type: "PlainText",
				Text: "Nothing",
			},
			ShouldEndSession: true,
		},
	}
}

func formattedDate(time time.Time) string {
	r := strings.NewReplacer(
		"January", "Enero",
		"February", "Febrero",
		"March", "Marzo",
		"April", "Abril",
		"May", "Mayo",
		"June", "Junio",
		"July", "Julio",
		"August", "Agosto",
		"September", "Septiembre",
		"October", "Octubre",
		"November", "Noviembre",
		"December", "Diciembre")
	return r.Replace(time.Format("2 de January"))
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
