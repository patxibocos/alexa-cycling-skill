package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
	"time"
)

const questionAttribute = "question"
const raceAttribute = "race"
const dayAttribute = "day"
const stageInfoAttributeValue = "StageInfo"
const raceSlot = "race"
const daySlot = "day"
const numberSlot = "number"

func tomorrow() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today.Add(24 * time.Hour)
}

func addStageInfoQuestionToSession(sessionAttributes map[string]interface{}, raceId string, day time.Time) {
	sessionAttributes[questionAttribute] = stageInfoAttributeValue
	sessionAttributes[raceAttribute] = raceId
	sessionAttributes[dayAttribute] = day.Format("2006-01-02")
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
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok && !rr.IsLastStage && cycling.StageContainsData(race.Stages[rr.StageNumber]) {
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceId, tomorrow())
		message += ". Quieres saber cómo es la etapa de mañana?"
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok && cycling.StageContainsData(race.Stages[0]) {
		if cycling.IsSingleDayRace(race) {
			message += ". Quieres saber cómo será la etapa?"
		} else {
			message += ". Quieres saber cómo será la primera etapa?"
		}
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceId, race.StartDate.AsTime())
	}
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
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
			if cycling.StageContainsData(nextRace.Stages[0]) {
				if cycling.IsSingleDayRace(nextRace) {
					message += ". Quieres saber cómo será la primera etapa?"
				} else {
					message += ". Quieres saber cómo será la etapa?"
				}
				endSession = false
				addStageInfoQuestionToSession(sessionAttributes, nextRace.Id, nextRace.StartDate.AsTime())
			}
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
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
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
	return newResponse().shouldEndSession(true).text(message)
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
	return newResponse().shouldEndSession(true).text(message)
}

func handleNo(_ Request, _ *pcsscraper.CyclingData) Response {
	return newResponse().shouldEndSession(true).text("")
}

func handleHelp(_ Request, _ *pcsscraper.CyclingData) Response {
	return newResponse().
		shouldEndSession(false).
		text("Me puedes preguntar por cuándo empieza o cómo va una carrera. También te puedo dar información sobre una etapa en particular. Qué quieres saber?")
}

func handleStop(_ Request, _ *pcsscraper.CyclingData) Response {
	return newResponse().shouldEndSession(true).text("¡Adios!")
}

func handleCancel(_ Request, _ *pcsscraper.CyclingData) Response {
	return newResponse().shouldEndSession(true).text("¡Adios!")
}

func handleYes(request Request, cyclingData *pcsscraper.CyclingData) Response {
	questionAttribute := request.Session.Attributes[questionAttribute]
	if questionAttribute == stageInfoAttributeValue {
		dayAttribute := request.Session.Attributes[dayAttribute]
		raceAttribute := request.Session.Attributes[raceAttribute]
		day, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", dayAttribute))
		raceId := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceId)
		raceStage := cycling.GetRaceStageForDay(race, day)
		message := messageForRaceStage(raceStage)
		return newResponse().shouldEndSession(true).text(message)
	}
	return newResponse().shouldEndSession(true).text("")
}
