package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
)

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
