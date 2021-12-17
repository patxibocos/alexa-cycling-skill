package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
)

func IntentDispatcher(request Request, cyclingData *pcsscraper.CyclingData) Response {
	message := ""
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
		raceName := RaceName(race.Id)
		switch ri := raceResult.(type) {
		case *cycling.PastRace:
			message = fmt.Sprintf(
				"%s terminó el %s. El ganador fue %s, el segundo %s y tercero %s",
				raceName,
				FormattedDate(race.EndDate.AsTime()),
				RiderFullName(ri.GcTop3.First),
				RiderFullName(ri.GcTop3.Second),
				RiderFullName(ri.GcTop3.Third),
			)
		case *cycling.FutureRace:
			message = fmt.Sprintf(
				"%s no empieza hasta el %s",
				raceName,
				FormattedDate(race.StartDate.AsTime()),
			)
		case *cycling.RestDayStage:
			message = fmt.Sprintf(
				"Hoy es día de descanso en %s",
				raceName,
			)
		case *cycling.SingleDayRaceWithResults:
			message = fmt.Sprintf(
				"Hoy se ha disputado %s. El ganador ha sido %s, el segundo %s y tercero %s",
				raceName,
				RiderFullName(ri.Top3.First),
				RiderFullName(ri.Top3.Second),
				RiderFullName(ri.Top3.Third),
			)
		case *cycling.SingleDayRaceWithoutResults:
			message = fmt.Sprintf(
				"Hoy se disputa %s pero todavía no tengo los resultados. Vuelve a preguntarme en un rato",
				raceName,
			)
		case *cycling.MultiStageRaceWithResults: // If stageNumber is greater than 1 -> return GC. If it is the last stage -> announce race has ended
			var stageName string
			if ri.IsLastStage {
				stageName = "última"
			} else {
				stageName = fmt.Sprintf("%dª", ri.StageNumber)
			}
			message = fmt.Sprintf(
				"Hoy se ha disputado la %s etapa de %s. El ganador ha sido %s, el segundo %s y tercero %s.",
				stageName,
				raceName,
				RiderFullName(ri.Top3.First),
				RiderFullName(ri.Top3.Second),
				RiderFullName(ri.Top3.Third),
			)
			if ri.StageNumber > 1 {
				message += fmt.Sprintf(
					"La clasificación general se queda: primero %s, segundo %s y tercero %s",
					RiderFullName(ri.GcTop3.First),
					RiderFullName(ri.GcTop3.Second),
					RiderFullName(ri.GcTop3.Third),
				)
			}
		case *cycling.MultiStageRaceWithoutResults:
			message = fmt.Sprintf(
				"Hoy se disputa la %dª etapa de %s pero todavía no tengo los resultados. Vuelve a preguntarme en un rato",
				ri.StageNumber,
				raceName,
			)
		}
	}
	return Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: message,
			},
			ShouldEndSession: true,
		},
	}
}
