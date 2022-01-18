package alexa

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
)

var intentRouting = map[string]func(Request, i18nLocalizer, *pcsscraper.CyclingData) Response{
	"RaceResult":          handleRaceResult,
	"DayStageInfo":        handleDayStageInfo,
	"NumberStageInfo":     handleNumberStageInfo,
	"MountainsStart":      handleMountainsStart,
	"AMAZON.YesIntent":    handleYes,
	"AMAZON.NoIntent":     handleNo,
	"AMAZON.HelpIntent":   handleHelp,
	"AMAZON.CancelIntent": handleCancel,
	"AMAZON.StopIntent":   handleStop,
}

func RequestHandler(request Request, cyclingData *pcsscraper.CyclingData) Response {
	localizer := newLocalizer(request.Body.Locale)
	if request.Body.Type == "LaunchRequest" {
		return handleLaunchRequest(request, localizer, cyclingData)
	}
	if request.Body.Type == "Connections.Response" {
		return handleConnectionsResponse(request, localizer, cyclingData)
	}
	if handler, ok := intentRouting[request.Body.Intent.Name]; ok {
		return handler(request, localizer, cyclingData)
	}
	return newResponse().shouldEndSession(true)
}
