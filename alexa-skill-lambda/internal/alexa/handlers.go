package alexa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"io"
	"io/ioutil"
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
const raceSlot = "race"
const daySlot = "day"
const numberSlot = "number"

var ErrUnauthorized = errors.New("ErrUnauthorized")

func today() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today
}

func tomorrow() time.Time {
	return today().Add(24 * time.Hour)
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
	var messages []string
	messages = append(messages, messageForRaceResult(localizer, race, raceResult))
	endSession := true
	sessionAttributes := make(map[string]interface{})
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok && !cycling.IsLastRaceStage(race, rr.StageNumber) && cycling.StageContainsData(race.Stages[rr.StageNumber]) {
		endSession = false
		addStageInfoQuestionToSession(sessionAttributes, raceId, tomorrow())
		messages = append(messages, localizer.localize(localizeParams{key: "TomorrowStageQuestion"}))
	}
	if _, ok := raceResult.(*cycling.FutureRace); ok {
		daysDiff := race.StartDate.AsTime().Sub(today()).Hours()
		if daysDiff >= 7 && !isReminderForRace(race, request) {
			messages = append(messages, localizer.localize(localizeParams{key: "RaceReminderQuestion"}))
			sessionAttributes[questionAttribute] = setReminderAttributeValue
			sessionAttributes[raceAttribute] = race.Id
			endSession = false
		} else if cycling.StageContainsData(race.Stages[0]) {
			if cycling.IsSingleDayRace(race) {
				messages = append(messages, localizer.localize(localizeParams{key: "SingleStageQuestion"}))
			} else {
				messages = append(messages, localizer.localize(localizeParams{key: "FistStageQuestion"}))
			}
			endSession = false
			addStageInfoQuestionToSession(sessionAttributes, raceId, race.StartDate.AsTime())
		}
	}
	message := strings.Join(messages, ". ")
	return newResponse().shouldEndSession(endSession).text(message).sessionAttributes(sessionAttributes)
}

func handleLaunchRequest(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
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
			daysDiff := nextRace.StartDate.AsTime().Sub(today()).Hours()
			if daysDiff >= 7 && !isReminderForRace(nextRace, request) {
				messages = append(messages, localizer.localize(localizeParams{key: "RaceReminderQuestion"}))
				sessionAttributes[questionAttribute] = setReminderAttributeValue
				sessionAttributes[raceAttribute] = nextRace.Id
				endSession = false
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

func handleConnectionsResponse(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	if request.Body.Payload.Status != "ACCEPTED" {
		return newResponse().text(localizer.localize(localizeParams{key: "Ooops"})).shouldEndSession(true)
	}
	splits := strings.Split(request.Body.Token, ":")
	action := splits[0]
	switch action {
	case setReminderAttributeValue:
		raceId := splits[1]
		race := cycling.FindRace(cyclingData.Races, raceId)
		err := setReminderForRace(request, localizer, race)
		if err != nil {
			return newResponse().text(localizer.localize(localizeParams{key: "SetReminderFailed"})).shouldEndSession(true)
		}
		message := localizer.localize(localizeParams{key: "RaceReminderSet"})
		return newResponse().text(message).shouldEndSession(true)
	}
	return newResponse().shouldEndSession(true)
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

func handleMountainsStart(request Request, localizer i18nLocalizer, cyclingData *pcsscraper.CyclingData) Response {
	intent := request.Body.Intent
	raceSlot := intent.Slots[raceSlot]
	raceId := raceSlot.Resolutions.ResolutionsPerAuthority[0].Values[0].Value.ID
	race := cycling.FindRace(cyclingData.Races, raceId)
	var message string
	mountainsStage := cycling.FindMountainsStage(race)
	raceName := raceName(raceId)
	switch ms := mountainsStage.(type) {
	case *cycling.SingleDayRace:
		message = localizer.localize(localizeParams{
			key: "MountainsSingleDayRace",
			data: map[string]interface{}{
				"Race":      raceName,
				"StartDate": formattedDate(race.StartDate.AsTime()),
			},
		})
	case *cycling.YesMountainsStage:
		message = localizer.localize(localizeParams{
			key: "MountainsStartAvailable",
			data: map[string]interface{}{
				"StartDate": formattedDate(ms.StartDate.AsTime()),
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

func handleNo(_ Request, _ i18nLocalizer, _ *pcsscraper.CyclingData) Response {
	return newResponse().shouldEndSession(true)
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
	switch questionAttribute {
	case stageInfoAttributeValue:
		dayAttribute := request.Session.Attributes[dayAttribute]
		raceAttribute := request.Session.Attributes[raceAttribute]
		day, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", dayAttribute))
		raceId := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceId)
		raceStage := cycling.GetRaceStageForDay(race, day)
		message := messageForRaceStage(localizer, raceStage)
		return newResponse().shouldEndSession(true).text(message)
	case setReminderAttributeValue:
		raceAttribute := request.Session.Attributes[raceAttribute]
		raceId := fmt.Sprintf("%v", raceAttribute)
		race := cycling.FindRace(cyclingData.Races, raceId)
		err := setReminderForRace(request, localizer, race)
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
				Token: fmt.Sprintf("%s:%s", setReminderAttributeValue, raceId),
			}
			return newResponse().directives([]interface{}{reminderDirective}).shouldEndSession(true)
		}
		if err != nil {
			return newResponse().text(localizer.localize(localizeParams{key: "SetReminderFailed"})).shouldEndSession(true)
		}
		message := localizer.localize(localizeParams{key: "RaceReminderSet"})
		return newResponse().text(message).shouldEndSession(true)
	}
	return newResponse().shouldEndSession(true)
}

func setReminderForRace(request Request, localizer i18nLocalizer, race *pcsscraper.Race) error {
	reminderMessage := localizer.localize(localizeParams{
		key: "RaceReminder",
		data: map[string]interface{}{
			"Race": raceName(race.Id),
		},
	})
	raceMillis := cycling.MillisForRace(race)
	reminderTime := race.StartDate.AsTime().Add(-14 * time.Hour).Add(time.Duration(raceMillis) * time.Millisecond)
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
	responseBytes, _ := ioutil.ReadAll(resp.Body)
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
