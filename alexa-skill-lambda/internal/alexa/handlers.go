package alexa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/timeutils"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const questionAttribute = "question"
const raceAttribute = "race"
const dayAttribute = "day"
const stageInfoAttributeValue = "StageInfo"
const setReminderAttributeValue = "SetReminder"
const raceGeneralClassificationAttributeValue = "RaceGeneralClassification"
const raceSlot = "race"
const daySlot = "day"
const numberSlot = "number"

var ErrUnauthorized = errors.New("ErrUnauthorized")

func addStageInfoQuestionToSession(sessionAttributes map[string]interface{}, raceId string, day time.Time) {
	sessionAttributes[questionAttribute] = stageInfoAttributeValue
	sessionAttributes[raceAttribute] = raceId
	sessionAttributes[dayAttribute] = day.Format("2006-01-02")
}

func handleRaceResult(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	location := locationProvider()
	intent := request.Body.Intent
	raceNameSlot := intent.Slots[raceSlot]
	raceIdPrefix := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
	raceResult := cycling.GetRaceResult(race, cyclingData.Riders, cyclingData.Teams, location)
	var messages []string
	messages = append(messages, messageForRaceResult(localizer, race, raceResult, location))
	endSession := true
	sessionAttributes := make(map[string]interface{})
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok && !cycling.IsLastRaceStage(race, rr.StageNumber) && cycling.StageContainsData(race.Stages[rr.StageNumber]) {
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceIdPrefix, timeutils.Tomorrow(location))
		messages = append(messages, localizer.localize(localizeParams{key: "TomorrowStageQuestion"}))
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok {
		hoursDiff := race.StartDateLocal(location).Sub(timeutils.Today(location)).Hours()
		if hoursDiff >= 24*7 && !isReminderForRace(race, request) {
			messages = append(messages, localizer.localize(localizeParams{key: "RaceReminderQuestion"}))
			sessionAttributes[questionAttribute] = setReminderAttributeValue
			sessionAttributes[raceAttribute] = race.Id
			endSession = false
		} else if cycling.StageContainsData(race.Stages[0]) {
			if cycling.IsSingleDayRace(race, location) {
				messages = append(messages, localizer.localize(localizeParams{key: "SingleStageQuestion"}))
			} else {
				messages = append(messages, localizer.localize(localizeParams{key: "FistStageQuestion"}))
			}
			endSession = false
			addStageInfoQuestionToSession(sessionAttributes, raceIdPrefix, race.StartDateLocal(location))
		}
	}
	message := strings.Join(messages, ". ")
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
}

func handleLaunchRequest(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	location := locationProvider()
	activeRaces := cycling.GetActiveRaces(cyclingData.Races, location)
	endSession := true
	var messages []string
	sessionAttributes := make(map[string]interface{})
	switch len(activeRaces) {
	case 0:
		messages = append(messages, localizer.localize(localizeParams{key: "NoActiveRace"}))
		nextRace := cycling.FindNextRace(cyclingData.Races, location)
		if nextRace == nil {
			messages = append(messages, localizer.localize(localizeParams{key: "SeasonEnded"}))
		} else {
			messages = append(messages, localizer.localize(localizeParams{
				key: "NextRaceStart",
				data: map[string]interface{}{
					"Race":      raceName(nextRace.Id),
					"StartDate": formattedDate(nextRace.StartDateLocal(location)),
				},
			}))
			hoursDiff := nextRace.StartDateLocal(location).Sub(timeutils.Today(location)).Hours()
			if hoursDiff >= 24*7 && !isReminderForRace(nextRace, request) {
				messages = append(messages, localizer.localize(localizeParams{key: "RaceReminderQuestion"}))
				sessionAttributes[questionAttribute] = setReminderAttributeValue
				sessionAttributes[raceAttribute] = nextRace.Id
				endSession = false
			}
		}
	case 1:
		race := activeRaces[0]
		raceResult := cycling.GetRaceResult(race, cyclingData.Riders, cyclingData.Teams, location)
		messages = append(messages, messageForRaceResult(localizer, race, raceResult, location))
		if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok && !cycling.IsLastRaceStage(race, rr.StageNumber) && cycling.StageContainsData(race.Stages[rr.StageNumber]) {
			endSession = false
			addStageInfoQuestionToSession(sessionAttributes, race.Id, timeutils.Tomorrow(location))
			messages = append(messages, localizer.localize(localizeParams{key: "TomorrowStageQuestion"}))
		}
		if rr, ok := raceResult.(*cycling.MultiStageRaceWithoutResults); ok && rr.StageNumber > 1 {
			endSession = false
			sessionAttributes[questionAttribute] = raceGeneralClassificationAttributeValue
			sessionAttributes[raceAttribute] = race.Id
			messages = append(messages, localizer.localize(localizeParams{key: "GeneralClassificationQuestion"}))
		}
		if _, ok := raceResult.(*cycling.RestDayStage); ok {
			endSession = false
			addStageInfoQuestionToSession(sessionAttributes, race.Id, timeutils.Tomorrow(location))
			messages = append(messages, localizer.localize(localizeParams{key: "TomorrowStageQuestion"}))
		}
	case 2:
		race1 := activeRaces[0]
		race2 := activeRaces[1]
		race1Result := cycling.GetRaceResult(race1, cyclingData.Riders, cyclingData.Teams, location)
		race2Result := cycling.GetRaceResult(race2, cyclingData.Riders, cyclingData.Teams, location)
		messages = append(messages, messageForTwoRaceResults(localizer, race1, race2, race1Result, race2Result, location))
	}
	message := strings.Join(messages, ". ")
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
}

func handleConnectionsResponse(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	if request.Body.Payload.Status != "ACCEPTED" {
		return newResponse().text(localizer.localize(localizeParams{key: "Ooops"})).shouldEndSession(true)
	}
	splits := strings.Split(request.Body.Token, ":")
	action := splits[0]
	switch action {
	case setReminderAttributeValue:
		location := locationProvider()
		raceIdPrefix := splits[1]
		race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
		err := setReminderForRace(request, localizer, race, location)
		if err != nil {
			return newResponse().text(localizer.localize(localizeParams{key: "SetReminderFailed"})).shouldEndSession(true)
		}
		message := localizer.localize(localizeParams{key: "RaceReminderSet"})
		return newResponse().text(message).shouldEndSession(true)
	}
	return newResponse().shouldEndSession(true)
}

func handleDayStageInfo(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	daySlot := intent.Slots[daySlot]
	day, _ := time.Parse("2006-01-02", daySlot.Value)
	raceIdPrefix := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
	raceStage := cycling.GetRaceStageForDay(race, day, locationProvider())
	message := messageForRaceStage(localizer, raceStage)
	return newResponse().shouldEndSession(true).text(message)
}

func handleNumberStageInfo(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	numberSlot := intent.Slots[numberSlot]
	raceIdPrefix := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
	stageIndex, _ := strconv.Atoi(numberSlot.Value)
	raceStage := cycling.GetRaceStageForIndex(race, stageIndex, locationProvider())
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
				"StartDate": formattedDate(rs.StartDateLocal),
			},
		})
	}
	return newResponse().shouldEndSession(true).text(message)
}

func handleMountainsStart(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	raceIdPrefix := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
	var message string
	location := locationProvider()
	mountainsStage := cycling.FindMountainsStage(race, location)
	raceName := raceName(raceIdPrefix)
	switch ms := mountainsStage.(type) {
	case *cycling.SingleDayRace:
		message = localizer.localize(localizeParams{
			key: "MountainsSingleDayRace",
			data: map[string]interface{}{
				"Race":      raceName,
				"StartDate": formattedDate(race.StartDateLocal(location)),
			},
		})
	case *cycling.YesMountainsStage:
		message = localizer.localize(localizeParams{
			key: "MountainsStartAvailable",
			data: map[string]interface{}{
				"StartDate": formattedDate(ms.StartDateLocal),
				"Race":      raceName,
			},
		})
	case *cycling.NoStageTypeData:
		message = localizer.localize(localizeParams{
			key: "MountainsNoTypeData",
			data: map[string]interface{}{
				"Race": raceName,
			},
		})
	case *cycling.NoMountainsStage:
		message = localizer.localize(localizeParams{
			key: "MountainsNoStage",
			data: map[string]interface{}{
				"Race": raceName,
			},
		})
	}
	return newResponse().shouldEndSession(true).text(message)
}

func handleGeneralClassification(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	intent := request.Body.Intent
	raceNameSlot := intent.Slots[raceSlot]
	raceIdPrefix := raceNameSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
	raceName := raceName(race.Id)
	if cycling.RaceHasNotStarted(race) {
		message := localizer.localize(localizeParams{
			key: "RaceResultFuture",
			data: map[string]interface{}{
				"Race":      raceName,
				"StartDate": formattedDate(race.StartDateLocal(locationProvider())),
			},
		})
		return newResponse().text(message).shouldEndSession(true)
	}
	gcTop3 := cycling.GetTop3FromResult(race.Result(), cyclingData.Riders)
	if cycling.RaceIsFinished(race) {
		message := localizer.localize(localizeParams{
			key: "GeneralClassificationComplete",
			data: map[string]interface{}{
				"First":  riderFullName(gcTop3.First.Rider),
				"Second": riderFullName(gcTop3.Second.Rider),
				"Third":  riderFullName(gcTop3.Third.Rider),
				"Race":   raceName,
			},
		})
		return newResponse().text(message).shouldEndSession(true)
	}
	message := localizer.localize(localizeParams{
		key: "RaceResultGeneralClassification",
		data: map[string]interface{}{
			"First":                riderFullName(gcTop3.First.Rider),
			"Second":               riderFullName(gcTop3.Second.Rider),
			"Third":                riderFullName(gcTop3.Third.Rider),
			"GapFromFirstToSecond": messageForGap(localizer, gcTop3.Second.Time-gcTop3.First.Time),
			"GapFromSecondToThird": messageForGap(localizer, gcTop3.Third.Time-gcTop3.Second.Time),
		},
	})
	return newResponse().text(message).shouldEndSession(true)
}

func handleNo(_ Request, _ i18nLocalizer, _ *pcsscraper.CyclingData, _ func() *time.Location) Response {
	return newResponse().shouldEndSession(true)
}

func handleHelp(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData, _ func() *time.Location) Response {
	message := localizer.localize(localizeParams{key: "Help"})
	return newResponse().
		shouldEndSession(false).
		text(message)
}

func handleStop(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData, _ func() *time.Location) Response {
	message := localizer.localize(localizeParams{key: "Goodbye"})
	return newResponse().shouldEndSession(true).text(message)
}

func handleCancel(_ Request, localizer i18nLocalizer, _ *pcsscraper.CyclingData, _ func() *time.Location) Response {
	message := localizer.localize(localizeParams{key: "Goodbye"})
	return newResponse().shouldEndSession(true).text(message)
}

func handleYes(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData, locationProvider func() *time.Location) Response {
	questionAttribute := request.Session.Attributes[questionAttribute]
	switch questionAttribute {
	case stageInfoAttributeValue:
		dayAttribute := request.Session.Attributes[dayAttribute]
		raceAttribute := request.Session.Attributes[raceAttribute]
		day, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", dayAttribute))
		raceIdPrefix := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
		raceStage := cycling.GetRaceStageForDay(race, day, locationProvider())
		message := messageForRaceStage(localizer, raceStage)
		return newResponse().shouldEndSession(true).text(message)
	case setReminderAttributeValue:
		raceAttribute := request.Session.Attributes[raceAttribute]
		raceIdPrefix := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
		err := setReminderForRace(request, localizer, race, locationProvider())
		if errors.Is(err, ErrUnauthorized) {
			reminderDirective := SendRequestDirective{
				Type: "Connections.SendRequest",
				Name: "AskFor",
				Payload: DirectivePayload{
					Type:    "AskForPermissionsConsentRequest",
					Version: "2",
					PermissionScopes: []DirectivePermissionScope{
						{
							PermissionScope: "alexa::alerts:reminders:skill:readwrite",
							ConsentLevel:    "ACCOUNT",
						},
					},
				},
				Token: fmt.Sprintf("%s:%s", setReminderAttributeValue, raceIdPrefix),
			}
			return newResponse().directives([]interface{}{reminderDirective}).shouldEndSession(true)
		}
		if err != nil {
			return newResponse().text(localizer.localize(localizeParams{key: "SetReminderFailed"})).shouldEndSession(true)
		}
		message := localizer.localize(localizeParams{key: "RaceReminderSet"})
		return newResponse().text(message).shouldEndSession(true)
	case raceGeneralClassificationAttributeValue:
		raceAttribute := request.Session.Attributes[raceAttribute]
		raceIdPrefix := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceIdPrefix)
		gcTop3 := cycling.GetTop3FromResult(race.Result(), cyclingData.Riders)
		message := localizer.localize(localizeParams{
			key: "RaceResultGeneralClassification",
			data: map[string]interface{}{
				"First":                riderFullName(gcTop3.First.Rider),
				"Second":               riderFullName(gcTop3.Second.Rider),
				"Third":                riderFullName(gcTop3.Third.Rider),
				"GapFromFirstToSecond": messageForGap(localizer, gcTop3.Second.Time-gcTop3.First.Time),
				"GapFromSecondToThird": messageForGap(localizer, gcTop3.Third.Time-gcTop3.Second.Time),
			},
		})
		return newResponse().text(message).shouldEndSession(true)
	}
	return newResponse().shouldEndSession(true)
}

func setReminderForRace(request Request, localizer i18nLocalizer, race *pcsscraper.Race, location *time.Location) error {
	reminderMessage := localizer.localize(localizeParams{
		key: "RaceReminder",
		data: map[string]interface{}{
			"Race": raceName(race.Id),
		},
	})
	raceMillis := cycling.MillisForRace(race)
	reminderTime := race.StartDateLocal(location).Add(-14 * time.Hour).Add(time.Duration(raceMillis) * time.Millisecond)
	reminderRequest := buildReminderRequest(reminderTime, request.Body.Locale, reminderMessage)
	serializedRequest, _ := json.Marshal(reminderRequest)
	resp, err := doRequest("POST", "/v1/alerts/reminders", request, bytes.NewBuffer(serializedRequest))
	if resp.StatusCode == 401 {
		return ErrUnauthorized
	}
	return err
}

func doRequest(method string, uri string, request Request, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", request.Context.System.ApiEndpoint, uri), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Context.System.ApiAccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func buildReminderRequest(scheduledTime time.Time, locale, text string) reminderRequest {
	return reminderRequest{
		RequestTime: time.Now().Format("2006-01-02T15:04:05.000"),
		Trigger: trigger{
			Type:          "SCHEDULED_ABSOLUTE",
			ScheduledTime: scheduledTime.Format("2006-01-02T15:04:05.000"),
		},
		AlertInfo: alertInfo{
			SpokenInfo: spokenInfo{
				Content: []content{{
					Locale: locale,
					Text:   text,
				}},
			},
		},
		PushNotification: pushNotification{
			Status: "ENABLED",
		},
	}
}

func isReminderForRace(race *pcsscraper.Race, request Request) bool {
	resp, _ := doRequest("GET", "/v1/alerts/reminders", request, nil)
	responseBytes, _ := io.ReadAll(resp.Body)
	response := new(remindersResponse)
	_ = json.Unmarshal(responseBytes, response)
	if response.TotalCount == "0" {
		return false
	}
	raceMillis := cycling.MillisForRace(race)
	for _, reminder := range response.Alerts {
		scheduledTime, _ := time.Parse("2006-01-02T15:04:05.000", reminder.Trigger.ScheduledTime)
		if int(scheduledTime.UnixMilli()%cycling.MillisModulo) == raceMillis {
			return true
		}
	}
	return false
}
