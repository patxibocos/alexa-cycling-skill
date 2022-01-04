package alexa

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
)

var intentRouting = map[string]func(Request, *pcsscraper.CyclingData) Response{
	"RaceResult":        handleRaceResult,
	"DayStageInfo":      handleDayStageInfo,
	"NumberStageInfo":   handleNumberStageInfo,
	"AMAZON.YesIntent":  handleYes,
	"AMAZON.NoIntent":   handleNo,
	"AMAZON.HelpIntent": handleHelp,
}

func RequestHandler(request Request, cyclingData *pcsscraper.CyclingData) Response {
	if request.Body.Type == "LaunchRequest" {
		return handleLaunchRequest(request, cyclingData)
	}
	if handler, ok := intentRouting[request.Body.Intent.Name]; ok {
		return handler(request, cyclingData)
	}
	return newResponse().shouldEndSession(true)
}
