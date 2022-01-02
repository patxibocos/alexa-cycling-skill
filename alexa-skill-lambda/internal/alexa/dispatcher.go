package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
	"time"
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

func handleDayStageInfo(intent Intent, cyclingData *pcsscraper.CyclingData) string {
	raceNameSlot := intent.Slots["raceName"]
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
	switch rs := raceStage.(type) {
	case *cycling.RestDayStage:
		return "Los corredores tienen descanso ese día"
	case *cycling.NoStage:
		return fmt.Sprintf("%s no tiene etapa para el %s", raceName(race.Id), formattedDate(day))
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
		return strings.Join(messages, ". ")
	case *cycling.StageWithoutData:
		return "Aún no hay información disponible de la etapa"
	}
	return ""
}

func handleNumberStageInfo(intent Intent, cyclingData *pcsscraper.CyclingData) string {
	raceNameSlot := intent.Slots["raceName"]
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
	switch rs := raceStage.(type) {
	case *cycling.NoStage:
		var stageOrStages string
		if len(race.Stages) > 1 {
			stageOrStages = "etapas"
		} else {
			stageOrStages = "etapa"
		}
		return fmt.Sprintf("%s sólo tiene %d %s", raceName(race.Id), len(race.Stages), stageOrStages)
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
		return strings.Join(messages, ". ")
	case *cycling.StageWithoutData:
		return fmt.Sprintf("La etapa comienza el %s. Aún no hay información disponible de la etapa", formattedDate(rs.StartDate.AsTime()))
	}
	return ""
}

func IntentDispatcher(request Request, cyclingData *pcsscraper.CyclingData) Response {
	message := ""
	if request.Body.Intent.Name == "RaceResult" {
		message = handleRaceResult(request.Body.Intent, cyclingData)
	} else if request.Body.Intent.Name == "DayStageInfo" {
		message = handleDayStageInfo(request.Body.Intent, cyclingData)
	} else if request.Body.Intent.Name == "NumberStageInfo" {
		message = handleNumberStageInfo(request.Body.Intent, cyclingData)
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
