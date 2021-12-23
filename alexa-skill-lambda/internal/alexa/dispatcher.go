package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strings"
)

func handleRaceResult(intent Intent, cyclingData *pcsscraper.CyclingData) string {
	raceNameSlot := intent.Slots["raceName"]
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	var race *pcsscraper.Race
	for _, r := range cyclingData.Races {
		if r.Id == raceId {
			race = r
		}
	}
	raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
	message := messageForRaceResult(race, raceResult)
	return message
}

func handleLaunchRequest(cyclingData *pcsscraper.CyclingData) string {
	activeRaces := cycling.GetActiveRaces(cyclingData.Races)
	switch len(activeRaces) {
	case 0:
		message := "No hay ninguna carrera activa ahora mismo"
		nextRace := cycling.FindNextRace(cyclingData.Races)
		if nextRace == nil {
			message += ". La temporada ha acabado"
		} else {
			message += fmt.Sprintf(". La siguiente carrera es %s y se disputa el %s", raceName(nextRace.Id), formattedDate(nextRace.StartDate.AsTime()))
		}
		return message
	case 1:
		race := activeRaces[0]
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
		return messageForRaceResult(race, raceResult)
	default:
		var raceMessages []string
		for _, race := range activeRaces {
			raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
			message := messageForRaceResult(race, raceResult)
			raceMessages = append(raceMessages, message)
		}
		return strings.Join(raceMessages, ". ")
	}
}

func IntentDispatcher(request Request, cyclingData *pcsscraper.CyclingData) Response {
	message := ""
	if request.Body.Intent.Name == "RaceResult" {
		message = handleRaceResult(request.Body.Intent, cyclingData)
	} else if request.Body.Type == "LaunchRequest" {
		message = handleLaunchRequest(cyclingData)
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

func messageForRaceResult(race *pcsscraper.Race, raceResult cycling.RaceResult) string {
	raceName := raceName(race.Id)
	switch ri := raceResult.(type) {
	case *cycling.PastRace:
		return fmt.Sprintf(
			"%s terminó el %s. %s",
			raceName,
			formattedDate(race.EndDate.AsTime()),
			phraseWithTop3("El ganador fue %s, segundo %s y tercero %s", ri.GcTop3),
		)
	case *cycling.FutureRace:
		return fmt.Sprintf(
			"%s no empieza hasta el %s",
			raceName,
			formattedDate(race.StartDate.AsTime()),
		)
	case *cycling.RestDayStage:
		return fmt.Sprintf(
			"Hoy es día de descanso en %s",
			raceName,
		)
	case *cycling.SingleDayRaceWithResults:
		return fmt.Sprintf(
			"Hoy se ha disputado %s. %s",
			raceName,
			phraseWithTop3("El ganador ha sido %s, segundo %s y tercero %s", ri.Top3),
		)
	case *cycling.SingleDayRaceWithoutResults:
		return fmt.Sprintf(
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
		message := fmt.Sprintf(
			"Hoy se ha disputado la %s etapa de %s. %s",
			stageName,
			raceName,
			phraseWithTop3("El ganador ha sido %s, segundo %s y tercero %s", ri.Top3),
		)
		if ri.StageNumber > 1 {
			message += phraseWithTop3AndGaps(". En la clasificación queda primero %s, segundo %s %s y tercero %s %s", ri.GcTop3)
		}
		return message
	case *cycling.MultiStageRaceWithoutResults:
		return fmt.Sprintf(
			"Hoy se disputa la %dª etapa de %s pero todavía no tengo los resultados. Vuelve a preguntarme en un rato",
			ri.StageNumber,
			raceName,
		)
	}
	return ""
}
