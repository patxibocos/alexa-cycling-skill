package alexa

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"time"
)

var intentRouting = map[string]func(Request, i18nLocalizer, *pcsscraper.CyclingData, func() *time.Location) Response{
	"RaceResult":            handleRaceResult,
	"DayStageInfo":          handleDayStageInfo,
	"NumberStageInfo":       handleNumberStageInfo,
	"MountainsStart":        handleMountainsStart,
	"GeneralClassification": handleGeneralClassification,
	"AMAZON.YesIntent":      handleYes,
	"AMAZON.NoIntent":       handleNo,
	"AMAZON.HelpIntent":     handleHelp,
	"AMAZON.CancelIntent":   handleCancel,
	"AMAZON.StopIntent":     handleStop,
}

func RequestHandler(request Request, cyclingData *pcsscraper.CyclingData) Response {
	localizer := newLocalizer(request.Body.Locale)
	if request.Body.Type == "LaunchRequest" {
		return handleLaunchRequest(request, localizer, cyclingData, locationProvider(request))
	}
	if request.Body.Type == "Connections.Response" {
		return handleConnectionsResponse(request, localizer, cyclingData, locationProvider(request))
	}
	if handler, ok := intentRouting[request.Body.Intent.Name]; ok {
		return handler(request, localizer, cyclingData, locationProvider(request))
	}
	return newResponse().shouldEndSession(true)
}

func locationProvider(request Request) func() *time.Location {
	return func() *time.Location {
		return getLocation(request)
	}
}
