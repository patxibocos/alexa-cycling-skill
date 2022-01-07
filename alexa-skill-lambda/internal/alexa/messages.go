package alexa

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
)

func messageForRaceStage(localizer i18nLocalizer, raceStage cycling.RaceStage) string {
	switch rs := raceStage.(type) {
	case *cycling.RestDayStage:
		return localizer.localize(localizeParams{key: "RaceStageRestDay"})
	case *cycling.NoStage:
		return localizer.localize(localizeParams{key: "RaceStageMissing"})
	case *cycling.StageWithData:
		return messageForStageWithData(localizer, rs)
	case *cycling.StageWithoutData:
		return localizer.localize(localizeParams{key: "RaceStageNoData"})
	}
	return ""
}

func messageForStageWithData(localizer i18nLocalizer, stageWithData *cycling.StageWithData) string {
	var messages []string
	if stageWithData.Departure != "" && stageWithData.Arrival != "" {
		message := localizer.localize(localizeParams{
			key:  "StageDataDepartureArrival",
			data: map[string]interface{}{"Departure": stageWithData.Departure, "Arrival": stageWithData.Arrival},
		})
		messages = append(messages, message)
	}
	if stageWithData.Distance > 0 {
		formattedDistance := strconv.FormatFloat(float64(stageWithData.Distance), 'f', -1, 32)
		message := localizer.localize(localizeParams{
			key:  "StageDataDistance",
			data: map[string]interface{}{"Distance": formattedDistance},
		})
		messages = append(messages, message)
	}
	if stageWithData.Type != pcsscraper.Stage_TYPE_UNSPECIFIED {
		message := localizer.localize(localizeParams{
			key:  "StageDataType",
			data: map[string]interface{}{"Type": messageForStageType(localizer, stageWithData.Type)},
		})
		messages = append(messages, message)
	}
	return strings.Join(messages, ". ")
}

func messageForRaceResult(localizer i18nLocalizer, race *pcsscraper.Race, raceResult cycling.RaceResult) string {
	raceName := raceName(race.Id)
	switch ri := raceResult.(type) {
	case *cycling.PastRace:
		return localizer.localize(localizeParams{
			key: "RaceResultPast",
			data: map[string]interface{}{
				"Race":    raceName,
				"EndDate": formattedDate(race.EndDate.AsTime()),
				"First":   riderFullName(ri.GcTop3.First.Rider),
				"Second":  riderFullName(ri.GcTop3.Second.Rider),
				"Third":   riderFullName(ri.GcTop3.Third.Rider),
			},
		})
	case *cycling.FutureRace:
		return localizer.localize(localizeParams{
			key: "RaceResultFuture",
			data: map[string]interface{}{
				"Race":      raceName,
				"StartDate": formattedDate(race.StartDate.AsTime()),
			},
		})
	case *cycling.RestDayStage:
		return localizer.localize(localizeParams{
			key: "RaceResultRestDay",
			data: map[string]interface{}{
				"Race": raceName,
			},
		})
	case *cycling.SingleDayRaceWithResults:
		return localizer.localize(localizeParams{
			key: "RaceResultSingleDayWithResults",
			data: map[string]interface{}{
				"Race":   raceName,
				"First":  riderFullName(ri.Top3.First.Rider),
				"Second": riderFullName(ri.Top3.Second.Rider),
				"Third":  riderFullName(ri.Top3.Third.Rider),
			},
		})
	case *cycling.SingleDayRaceWithoutResults:
		return localizer.localize(localizeParams{
			key: "RaceResultSingleDayWithoutResults",
			data: map[string]interface{}{
				"Race": raceName,
			},
		})
	case *cycling.MultiStageRaceWithResults: // If stageNumber is greater than 1 -> return GC. If it is the last stage -> announce race has ended
		stageName := messageForStageName(localizer, race, ri.StageNumber)
		message := localizer.localize(localizeParams{
			key: "RaceResultMultiStageWithResults",
			data: map[string]interface{}{
				"StageName": stageName,
				"Race":      raceName,
				"First":     riderFullName(ri.Top3.First.Rider),
				"Second":    riderFullName(ri.Top3.Second.Rider),
				"Third":     riderFullName(ri.Top3.Third.Rider),
			},
		})
		if ri.StageNumber > 1 {
			message += localizer.localize(localizeParams{
				key: "RaceResultGeneralClassification",
				data: map[string]interface{}{
					"First":                riderFullName(ri.GcTop3.First.Rider),
					"Second":               riderFullName(ri.GcTop3.Second.Rider),
					"Third":                riderFullName(ri.GcTop3.Third.Rider),
					"GapFromFirstToSecond": messageForGap(localizer, ri.GcTop3.Second.Time-ri.GcTop3.First.Time),
					"GapFromSecondToThird": messageForGap(localizer, ri.GcTop3.Third.Time-ri.GcTop3.Second.Time),
				},
			})
		}
		return message
	case *cycling.MultiStageRaceWithoutResults:
		stageName := messageForStageName(localizer, race, ri.StageNumber)
		return localizer.localize(localizeParams{
			key: "RaceResultMultiStageWithoutResults",
			data: map[string]interface{}{
				"StageName": stageName,
				"Race":      raceName,
			},
		})
	}
	return ""
}

func messageForStageName(localizer i18nLocalizer, race *pcsscraper.Race, stageNumber int) string {
	var stageName string
	if stageNumber == 1 {
		stageName = localizer.localize(localizeParams{key: "FirstStage"})
	} else if cycling.IsLastRaceStage(race, stageNumber) {
		stageName = localizer.localize(localizeParams{key: "LastStage"})
	} else {
		stageName = localizer.localize(localizeParams{key: "NumberStage", data: map[string]interface{}{"Number": stageNumber}})
	}
	return stageName
}

func messageForStageType(localizer i18nLocalizer, stageType pcsscraper.Stage_Type) string {
	var messageKey string
	switch stageType {
	case pcsscraper.Stage_TYPE_FLAT:
		messageKey = "StageTypeFlat"
	case pcsscraper.Stage_TYPE_HILLS_FLAT_FINISH:
		messageKey = "StageTypeHillsFlatFinish"
	case pcsscraper.Stage_TYPE_HILLS_UPHILL_FINISH:
		messageKey = "StageTypeHillsUphillFinish"
	case pcsscraper.Stage_TYPE_MOUNTAINS_FLAT_FINISH:
		messageKey = "StageTypeMountainsFlatFinish"
	case pcsscraper.Stage_TYPE_MOUNTAINS_UPHILL_FINISH:
		messageKey = "StageTypeMountainsUphillFinish"
	}
	return localizer.localize(localizeParams{key: messageKey})
}

func messageForGap(localizer i18nLocalizer, gap int64) string {
	if gap == 0 {
		return localizer.localize(localizeParams{key: "GapSameTime"})
	}
	if gap < 60 {
		return localizer.localize(localizeParams{
			key: "GapSeconds",
			data: map[string]interface{}{
				"Seconds": gap,
			},
			pluralCount: gap,
		})
	}
	minutes := gap / 60
	seconds := gap % 60
	message := localizer.localize(localizeParams{
		key: "GapMinutes",
		data: map[string]interface{}{
			"Minutes": minutes,
		},
		pluralCount: minutes,
	})
	if seconds > 0 {
		message += localizer.localize(localizeParams{
			key: "GapAndSeconds",
			data: map[string]interface{}{
				"Seconds": seconds,
			},
			pluralCount: seconds,
		})
	}
	return message
}
