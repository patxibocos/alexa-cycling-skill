package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
	"time"
)

func handleRaceResult(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots["race"]
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	var race *pcsscraper.Race
	for _, r := range cyclingData.Races {
		if r.Id == raceId {
			race = r
		}
	}
	raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
	message := messageForRaceResult(race, raceResult)
	endSession := true
	sessionAttributes := make(map[string]interface{})
	if multiStageRaceWithResults, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok {
		if !multiStageRaceWithResults.IsLastStage {
			// Ask if user wants info about next race
			// Set in session raceId and tomorrow's date to handle in YesIntent handler
			endSession = false
			sessionAttributes["question"] = "StageInfo"
			sessionAttributes["race"] = raceId
			sessionAttributes["day"] = cycling.Today().Add(24 * time.Hour).Format("2006-01-02")
			message += ". Quieres saber cómo es la etapa de mañana?"
		}
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok {
		endSession = false
		sessionAttributes["question"] = "StageInfo"
		sessionAttributes["race"] = raceId
		sessionAttributes["day"] = race.StartDate.AsTime().Format("2006-01-02")
		message += ". Quieres saber cómo será la primera etapa?"
	}
	return Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: message,
			},
			ShouldEndSession: endSession,
		},
	}
}

func handleLaunchRequest(_ Request, cyclingData *pcsscraper.CyclingData) Response {
	activeRaces := cycling.GetActiveRaces(cyclingData.Races)
	endSession := true
	var message string
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
			message += ". Quieres saber cómo será la primera etapa?"
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
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: message,
			},
			ShouldEndSession: endSession,
		},
	}
}

func handleDayStageInfo(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots["race"]
	daySlot := intent.Slots["day"]
	day, _ := time.Parse("2006-01-02", daySlot.Value)
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	var race *pcsscraper.Race
	for _, r := range cyclingData.Races {
		if r.Id == raceId {
			race = r
		}
	}
	raceStage := cycling.GetRaceStageForDay(race, day)
	var message string
	switch rs := raceStage.(type) {
	case *cycling.RestDayStage:
		message = "Los corredores tienen descanso ese día"
	case *cycling.NoStage:
		message = fmt.Sprintf("%s no tiene etapa para el %s", raceName(race.Id), formattedDate(day))
	case *cycling.StageWithData:
		var messages []string
		if rs.Departure != "" && rs.Arrival != "" {
			messages = append(messages, fmt.Sprintf("El recorrido va de %s a %s", rs.Departure, rs.Arrival))
		}
		if rs.Distance > 0 {
			formattedDistance := strconv.FormatFloat(float64(rs.Distance), 'f', -1, 32)
			messages = append(messages, fmt.Sprintf("Tiene una distancia de %s kilómetros", formattedDistance))
		}
		if rs.Type != pcsscraper.Stage_TYPE_UNSPECIFIED {
			messages = append(messages, fmt.Sprintf("El perfil de la etapa es %s", stageType(rs.Type)))
		}
		message = strings.Join(messages, ". ")
	case *cycling.StageWithoutData:
		message = "Aún no hay información disponible de la etapa"
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

func handleNumberStageInfo(request Request, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots["race"]
	numberSlot := intent.Slots["number"]
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	var race *pcsscraper.Race
	for _, r := range cyclingData.Races {
		if r.Id == raceId {
			race = r
		}
	}
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
		var messages []string
		messages = append(messages, fmt.Sprintf("La etapa comienza el %s", formattedDate(rs.StartDate.AsTime())))
		if rs.Departure != "" && rs.Arrival != "" {
			messages = append(messages, fmt.Sprintf("El recorrido va de %s a %s", rs.Departure, rs.Arrival))
		}
		if rs.Distance > 0 {
			formattedDistance := strconv.FormatFloat(float64(rs.Distance), 'f', -1, 32)
			messages = append(messages, fmt.Sprintf("Tiene una distancia de %s kilómetros", formattedDistance))
		}
		if rs.Type != pcsscraper.Stage_TYPE_UNSPECIFIED {
			messages = append(messages, fmt.Sprintf("El perfil de la etapa es %s", stageType(rs.Type)))
		}
		message = strings.Join(messages, ". ")
	case *cycling.StageWithoutData:
		message = fmt.Sprintf("La etapa comienza el %s. Aún no hay información disponible de la etapa", formattedDate(rs.StartDate.AsTime()))
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

func handleNo(_ Request, _ *pcsscraper.CyclingData) Response {
	return Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: "Sí",
			},
			ShouldEndSession: true,
		},
	}
}

func handleYes(_ Request, _ *pcsscraper.CyclingData) Response {
	return Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: "No",
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
