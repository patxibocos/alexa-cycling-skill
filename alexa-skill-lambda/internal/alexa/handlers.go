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

func handleRaceResult(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots[raceSlot]
	raceId := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
	message := messageForRaceResult(localizer, race, raceResult)
	endSession := true
	sessionAttributes := make(map[string]interface{})
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok && !cycling.IsLastRaceStage(race, rr.StageNumber) && cycling.StageContainsData(race.Stages[rr.StageNumber]) {
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceId, tomorrow())
		message += ". Quieres saber cómo es la etapa de mañana?" // TODO localizer.localize("TomorrowStageQuestion")
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok && cycling.StageContainsData(race.Stages[0]) {
		if cycling.IsSingleDayRace(race) {
			message += ". Quieres saber cómo será la etapa?" // TODO localizer.localize("SingleStageQuestion")
		} else {
			message += ". Quieres saber cómo será la primera etapa?" // TODO localizer.localize("FistStageQuestion")
		}
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceId, race.StartDate.AsTime())
	}
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
}

func handleLaunchRequest(_ Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	activeRaces := cycling.GetActiveRaces(cyclingData.Races)
	endSession := true
	var messages []string
	sessionAttributes := make(map[string]interface{})
	switch len(activeRaces) {
	case 0:
		messages = append(messages, localizer.localize(localizeParams{key: "NoActiveRace"}))
		nextRace := cycling.FindNextRace(cyclingData.Races)
		if nextRace == nil {
			messages = append(messages, localizer.localize(localizeParams{key: "SeasonEnded"}))
		} else {
			messages = append(messages, localizer.localize(localizeParams{
				key: "NextRaceStart",
				data: map[string]interface{}{
					"Race":      raceName(nextRace.Id),
					"StartDate": formattedDate(nextRace.StartDate.AsTime()),
				},
			}))
			if cycling.StageContainsData(nextRace.Stages[0]) {
				if cycling.IsSingleDayRace(nextRace) {
					messages = append(messages, localizer.localize(localizeParams{key: "FistStageQuestion"}))
				} else {
					messages = append(messages, localizer.localize(localizeParams{key: "SingleStageQuestion"}))
				}
				endSession = false
				addStageInfoQuestionToSession(sessionAttributes, nextRace.Id, nextRace.StartDate.AsTime())
			}
		}
	case 1:
		race := activeRaces[0]
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
		messages = append(messages, messageForRaceResult(localizer, race, raceResult))
	default:
		for _, race := range activeRaces {
			raceResult := cycling.GetRaceResult(race, cyclingData.Riders)
			message := messageForRaceResult(localizer, race, raceResult)
			messages = append(messages, message)
		}
	}
	message := strings.Join(messages, ". ")
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
}

func handleDayStageInfo(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	daySlot := intent.Slots[daySlot]
	day, _ := time.Parse("2006-01-02", daySlot.Value)
	raceId := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	raceStage := cycling.GetRaceStageForDay(race, day)
	message := messageForRaceStage(localizer, raceStage)
	return newResponse().shouldEndSession(true).text(message)
}

func handleNumberStageInfo(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
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
		message = localizer.localize(localizeParams{
			key: "RaceStages",
			data: map[string]interface{}{
				"Race":   raceName(race.Id),
				"Stages": len(race.Stages),
			},
			pluralCount: len(race.Stages),
		})
	case *cycling.StageWithData:
		message = messageForStageWithData(localizer, rs)
	case *cycling.StageWithoutData:
		message = localizer.localize(localizeParams{
			key: "NoDataForStage",
			data: map[string]interface{}{
				"StartDate": formattedDate(rs.StartDate.AsTime()),
			},
		})
	}
	return newResponse().shouldEndSession(true).text(message)
}

func handleNo(_ Request, _ i18nLocalizer, _ *pcsscraper.CyclingData) Response {
	return newResponse().shouldEndSession(true).text("")
}

func handleHelp(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData) Response {
	message := localizer.localize(localizeParams{key: "Help"})
	return newResponse().
		shouldEndSession(false).
		text(message)
}

func handleStop(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData) Response {
	message := localizer.localize(localizeParams{key: "Goodbye"})
	return newResponse().shouldEndSession(true).text(message)
}

func handleCancel(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData) Response {
	message := localizer.localize(localizeParams{key: "Goodbye"})
	return newResponse().shouldEndSession(true).text(message)
}

func handleYes(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	questionAttribute := request.Session.Attributes[questionAttribute]
	if questionAttribute == stageInfoAttributeValue {
		dayAttribute := request.Session.Attributes[dayAttribute]
		raceAttribute := request.Session.Attributes[raceAttribute]
		day, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", dayAttribute))
		raceId := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceId)
		raceStage := cycling.GetRaceStageForDay(race, day)
		message := messageForRaceStage(localizer, raceStage)
		return newResponse().shouldEndSession(true).text(message)
	}
	return newResponse().shouldEndSession(true).text("")
}
