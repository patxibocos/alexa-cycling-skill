package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
	"time"
)

const version = "1.0"
const plainText = "PlainText"
const questionAttribute = "question"
const raceAttribute = "race"
const dayAttribute = "day"
const raceSlot = "race"
const daySlot = "day"
const numberSlot = "number"

func tomorrow() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today.Add(24 * time.Hour)
}

func handleRaceResult(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots[raceSlot]
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
	message := messageForRaceResult(race, raceResult)
	endSession := true
	sessionAttributes := make(map[string]interface{})
	if multiStageRaceWithResults, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok {
		if !multiStageRaceWithResults.IsLastStage {
			endSession = false
			sessionAttributes[questionAttribute] = "StageInfo"
			sessionAttributes[raceAttribute] = raceId
			sessionAttributes[dayAttribute] = tomorrow().Format("2006-01-02")
			message += ". Quieres saber cómo es la etapa de mañana?"
		}
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok {
		// TODO Only propose getting info about next stage in case there is info available
		endSession = false
		sessionAttributes[questionAttribute] = "StageInfo"
		sessionAttributes[raceAttribute] = raceId
		sessionAttributes[dayAttribute] = race.StartDate.AsTime().Format("2006-01-02")
		if len(race.Stages) > 1 {
			message += ". Quieres saber cómo será la primera etapa?"
		} else {
			message += ". Quieres saber cómo será la etapa?"
		}
	}
	return Response{
		Version: version,
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
				Text: message,
			},
			ShouldEndSession: endSession,
		},
		SessionAttributes: sessionAttributes,
	}
}

func handleLaunchRequest(_ Request, cyclingData *pcsscraper.CyclingData) Response {
	activeRaces := cycling.GetActiveRaces(cyclingData.Races)
	endSession := true
	var message string
	sessionAttributes := make(map[string]interface{})
	switch len(activeRaces) {
	case 0:
		message = "No hay ninguna carrera activa ahora mismo"
		nextRace := cycling.FindNextRace(cyclingData.Races)
		if nextRace == nil {
			message += ". La temporada ha acabado"
		} else {
			message += fmt.Sprintf(
				". La siguiente carrera es %s y se disputa el %s",
				raceName(nextRace.Id),
				formattedDate(nextRace.StartDate.AsTime()),
			)
			if len(nextRace.Stages) > 1 {
				message += ". Quieres saber cómo será la primera etapa?"
			} else {
				message += ". Quieres saber cómo será la etapa?"
			}
			sessionAttributes[questionAttribute] = "StageInfo"
			sessionAttributes[raceAttribute] = nextRace.Id
			sessionAttributes[dayAttribute] = nextRace.StartDate.AsTime().Format("2006-01-02")
			endSession = false
		}
	case 1:
		race := activeRaces[0]
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
		message = messageForRaceResult(race, raceResult)
	default:
		var raceMessages []string
		for _, race := range activeRaces {
			raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
			message := messageForRaceResult(race, raceResult)
			raceMessages = append(raceMessages, message)
		}
		message = strings.Join(raceMessages, ". ")
	}
	return Response{
		Version:           version,
		SessionAttributes: sessionAttributes,
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
				Text: message,
			},
			ShouldEndSession: endSession,
		},
	}
}

func messageForRaceStage(raceStage cycling.RaceStage) string {
	var message string
	switch rs := raceStage.(type) {
	case *cycling.RestDayStage:
		message = "Los corredores tienen descanso"
	case *cycling.NoStage:
		message = "No hay etapa para ese día"
	case *cycling.StageWithData:
		message = messageForStageWithData(rs)
	case *cycling.StageWithoutData:
		message = "Aún no hay información disponible de la etapa"
	}
	return message
}

func messageForStageWithData(stageWithData *cycling.StageWithData) string {
	var messages []string
	if stageWithData.Departure != "" && stageWithData.Arrival != "" {
		messages = append(messages, fmt.Sprintf("El recorrido va de %s a %s", stageWithData.Departure, stageWithData.Arrival))
	}
	if stageWithData.Distance > 0 {
		formattedDistance := strconv.FormatFloat(float64(stageWithData.Distance), 'f', -1, 32)
		messages = append(messages, fmt.Sprintf("Tiene una distancia de %s kilómetros", formattedDistance))
	}
	if stageWithData.Type != pcsscraper.Stage_TYPE_UNSPECIFIED {
		messages = append(messages, fmt.Sprintf("El perfil de la etapa es %s", stageType(stageWithData.Type)))
	}
	return strings.Join(messages, ". ")
}

func handleDayStageInfo(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	daySlot := intent.Slots[daySlot]
	day, _ := time.Parse("2006-01-02", daySlot.Value)
	raceId := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	raceStage := cycling.GetRaceStageForDay(race, day)
	message := messageForRaceStage(raceStage)
	return Response{
		Version: version,
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
				Text: message,
			},
			ShouldEndSession: true,
		},
	}
}

func handleNumberStageInfo(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	numberSlot := intent.Slots[numberSlot]
	raceId := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	stageIndex, _ := strconv.Atoi(numberSlot.Value)
	raceStage := cycling.GetRaceStageForIndex(race, stageIndex)
	var message string
	switch rs := raceStage.(type) {
	case *cycling.NoStage:
		var stageOrStages string
		if len(race.Stages) > 1 {
			stageOrStages = "etapas"
		} else {
			stageOrStages = "etapa"
		}
		message = fmt.Sprintf("%s sólo tiene %d %s", raceName(race.Id), len(race.Stages), stageOrStages)
	case *cycling.StageWithData:
		message = messageForStageWithData(rs)
	case *cycling.StageWithoutData:
		message = fmt.Sprintf("La etapa comienza el %s pero aún no hay información disponible", formattedDate(rs.StartDate.AsTime()))
	}
	return Response{
		Version: version,
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
				Text: message,
			},
			ShouldEndSession: true,
		},
	}
}

func handleNo(_ Request, _ *pcsscraper.CyclingData) Response {
	return Response{
		Version:           version,
		SessionAttributes: make(map[string]interface{}),
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
				Text: "",
			},
			ShouldEndSession: true,
		},
	}
}

func handleYes(request Request, cyclingData *pcsscraper.CyclingData) Response {
	// TODO Check that question is "StageInfo"
	dayAttribute := request.Session.Attributes[dayAttribute]
	raceAttribute := request.Session.Attributes[raceAttribute]
	day, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", dayAttribute))
	raceId := fmt.Sprintf("%v", raceAttribute)
	race := cycling.FindRace(cyclingData.Races, raceId)
	raceStage := cycling.GetRaceStageForDay(race, day)
	message := messageForRaceStage(raceStage)
	return Response{
		Version:           version,
		SessionAttributes: make(map[string]interface{}),
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: plainText,
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
