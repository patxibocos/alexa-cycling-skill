package alexa

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strconv"
	"strings"
	"time"
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
	if stageWithData.ProfileType != pcsscraper.Stage_PROFILE_TYPE_UNSPECIFIED {
		message := localizer.localize(localizeParams{
			key:  "StageDataType",
			data: map[string]interface{}{"Type": messageForStageType(localizer, stageWithData.ProfileType)},
		})
		messages = append(messages, message)
	}
	return strings.Join(messages, ". ")
}

func messageForTwoRaceResults(localizer i18nLocalizer, race1, race2 *pcsscraper.Race, race1Result, race2Result cycling.RaceResult, location *time.Location) string {
	race1Name := messageForRaceOrStage(localizer, race1, race1Result, location)
	race2Name := messageForRaceOrStage(localizer, race2, race2Result, location)
	// No results
	if !raceResultsAvailable(race1Result) && !raceResultsAvailable(race2Result) {
		return localizer.localize(localizeParams{
			key:  "TwoRacesWithoutResults",
			data: map[string]interface{}{"Race1": race1Name, "Race2": race2Name},
		})
	}
	if raceResultsAvailable(race1Result) && raceResultsAvailable(race2Result) {
		var messages []string
		messages = append(messages, localizer.localize(localizeParams{
			key:  "TwoRacesWithResults",
			data: map[string]interface{}{"FullRaceName1": race1Name, "FullRaceName2": race2Name},
		}))
		messages = append(messages, messageForOneRaceWhenBothResultsAvailable(localizer, race1, race1Result, location))
		messages = append(messages, messageForOneRaceWhenBothResultsAvailable(localizer, race2, race2Result, location))
		return strings.Join(messages, ". ")
	}
	if raceResultsAvailable(race1Result) {
		return messageForTwoRacesWithResultForOne(localizer, race1, race2, race1Result, race2Result, location)
	}
	return messageForTwoRacesWithResultForOne(localizer, race2, race1, race2Result, race1Result, location)
}

func messageForOneRaceWhenBothResultsAvailable(localizer i18nLocalizer, race *pcsscraper.Race, raceResult cycling.RaceResult, location *time.Location) string {
	raceName := raceName(race.Id)
	if cycling.IsSingleDayRace(race, location) {
		rr, _ := raceResult.(*cycling.SingleDayRaceWithResults)
		return localizer.localize(localizeParams{
			key: "TwoRacesWithResultsSingleDayRace",
			data: map[string]interface{}{
				"Race":   raceName,
				"First":  riderFullName(rr.Top3.First.Rider),
				"Second": riderFullName(rr.Top3.Second.Rider),
				"Third":  riderFullName(rr.Top3.Third.Rider),
			},
		})
	}
	rr, _ := raceResult.(*cycling.MultiStageRaceWithResults)
	return localizer.localize(localizeParams{
		key: "TwoRacesWithResultsMultiStageRace",
		data: map[string]interface{}{
			"Race":   raceName,
			"Winner": riderFullName(rr.Top3.First.Rider),
			"First":  riderFullName(rr.GcTop3.First.Rider),
			"Second": riderFullName(rr.GcTop3.Second.Rider),
			"Third":  riderFullName(rr.GcTop3.Third.Rider),
		},
	})
}

func messageForTwoRacesWithResultForOne(localizer i18nLocalizer, raceWithResult, raceWithoutResult *pcsscraper.Race, resultForRaceWithResult, resultForRaceWithoutResult cycling.RaceResult, location *time.Location) string {
	raceWithResultName := messageForRaceOrStage(localizer, raceWithResult, resultForRaceWithResult, location)
	raceWithoutResultName := messageForRaceOrStage(localizer, raceWithoutResult, resultForRaceWithoutResult, location)
	var messages []string
	if cycling.IsSingleDayRace(raceWithResult, location) {
		rr, _ := resultForRaceWithResult.(*cycling.SingleDayRaceWithResults)
		messages = append(messages, localizer.localize(localizeParams{
			key: "RaceResultSingleDayWithResults",
			data: map[string]interface{}{
				"Race":   raceWithResultName,
				"First":  riderFullName(rr.Top3.First.Rider),
				"Second": riderFullName(rr.Top3.Second.Rider),
				"Third":  riderFullName(rr.Top3.Third.Rider),
			},
		}))
	} else {
		rr, _ := resultForRaceWithResult.(*cycling.MultiStageRaceWithResults)
		messages = append(messages, localizer.localize(localizeParams{
			key: "RaceResultMultiStageWithResults",
			data: map[string]interface{}{
				"MultiStageRaceName": raceWithResultName,
				"First":              riderFullName(rr.Top3.First.Rider),
				"Second":             riderFullName(rr.Top3.Second.Rider),
				"Third":              riderFullName(rr.Top3.Third.Rider),
			},
		}))
		if rr.StageNumber > 1 {
			messages = append(messages, localizer.localize(localizeParams{
				key: "RaceResultGeneralClassification",
				data: map[string]interface{}{
					"First":                riderFullName(rr.GcTop3.First.Rider),
					"Second":               riderFullName(rr.GcTop3.Second.Rider),
					"Third":                riderFullName(rr.GcTop3.Third.Rider),
					"GapFromFirstToSecond": messageForGap(localizer, rr.GcTop3.Second.Time-rr.GcTop3.First.Time),
					"GapFromSecondToThird": messageForGap(localizer, rr.GcTop3.Third.Time-rr.GcTop3.Second.Time),
				},
			}))
		}
	}
	messages = append(messages, localizer.localize(localizeParams{
		key: "TwoRacesWithSingleResult",
		data: map[string]interface{}{
			"Race": raceWithoutResultName,
		},
	}))
	return strings.Join(messages, ". ")
}

func raceResultsAvailable(raceResult cycling.RaceResult) bool {
	_, singleDayResults := raceResult.(*cycling.SingleDayRaceWithResults)
	_, multiStageResults := raceResult.(*cycling.MultiStageRaceWithResults)
	return singleDayResults || multiStageResults
}

func messageForRaceResult(localizer i18nLocalizer, race *pcsscraper.Race, raceResult cycling.RaceResult, location *time.Location) string {
	raceName := raceName(race.Id)
	switch ri := raceResult.(type) {
	case *cycling.PastRace:
		return localizer.localize(localizeParams{
			key: "RaceResultPast",
			data: map[string]interface{}{
				"Race":    raceName,
				"EndDate": formattedDate(race.EndDateLocal(location)),
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
				"StartDate": formattedDate(race.StartDateLocal(location)),
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
		multiStageRaceName := messageForRaceOrStage(localizer, race, ri, location)
		var messages []string
		messages = append(messages, localizer.localize(localizeParams{
			key: "RaceResultMultiStageWithResults",
			data: map[string]interface{}{
				"MultiStageRaceName": multiStageRaceName,
				"First":              riderFullName(ri.Top3.First.Rider),
				"Second":             riderFullName(ri.Top3.Second.Rider),
				"Third":              riderFullName(ri.Top3.Third.Rider),
			},
		}))
		if ri.StageNumber > 1 {
			messages = append(messages, localizer.localize(localizeParams{
				key: "RaceResultGeneralClassification",
				data: map[string]interface{}{
					"First":                riderFullName(ri.GcTop3.First.Rider),
					"Second":               riderFullName(ri.GcTop3.Second.Rider),
					"Third":                riderFullName(ri.GcTop3.Third.Rider),
					"GapFromFirstToSecond": messageForGap(localizer, ri.GcTop3.Second.Time-ri.GcTop3.First.Time),
					"GapFromSecondToThird": messageForGap(localizer, ri.GcTop3.Third.Time-ri.GcTop3.Second.Time),
				},
			}))
		}
		return strings.Join(messages, ". ")
	case *cycling.MultiTTTStageRaceWithResults:
		multiStageRaceName := messageForRaceOrStage(localizer, race, ri, location)
		return localizer.localize(localizeParams{
			key: "RaceResultMultiTTTStageWithResults",
			data: map[string]interface{}{
				"MultiStageRaceName": multiStageRaceName,
				"First":              ri.Top3Teams.First.Team.Name,
				"Second":             ri.Top3Teams.Second.Team.Name,
				"Third":              ri.Top3Teams.Third.Team.Name,
			},
		})
	case *cycling.MultiStageRaceWithoutResults:
		multiStageRaceName := messageForRaceOrStage(localizer, race, ri, location)
		return localizer.localize(localizeParams{
			key: "RaceResultMultiStageWithoutResults",
			data: map[string]interface{}{
				"MultiStageRaceName": multiStageRaceName,
			},
		})
	}
	return ""
}

func messageForRaceOrStage(localizer i18nLocalizer, race *pcsscraper.Race, raceResult cycling.RaceResult, location *time.Location) string {
	raceName := raceName(race.Id)
	if cycling.IsSingleDayRace(race, location) {
		return raceName
	}
	stageName := ""
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithoutResults); ok {
		stageName = messageForStageName(localizer, race, rr.StageNumber)
	}
	if rr, ok := raceResult.(*cycling.MultiStageRaceWithResults); ok {
		stageName = messageForStageName(localizer, race, rr.StageNumber)
	}
	if rr, ok := raceResult.(*cycling.MultiTTTStageRaceWithResults); ok {
		stageName = messageForStageName(localizer, race, rr.StageNumber)
	}
	return localizer.localize(localizeParams{
		key: "MultiStageRaceName",
		data: map[string]interface{}{
			"StageName": stageName,
			"Race":      raceName,
		},
	})
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

func messageForStageType(localizer i18nLocalizer, stageType pcsscraper.Stage_ProfileType) string {
	var messageKey string
	switch stageType {
	case pcsscraper.Stage_PROFILE_TYPE_FLAT:
		messageKey = "StageTypeFlat"
	case pcsscraper.Stage_PROFILE_TYPE_HILLS_FLAT_FINISH:
		messageKey = "StageTypeHillsFlatFinish"
	case pcsscraper.Stage_PROFILE_TYPE_HILLS_UPHILL_FINISH:
		messageKey = "StageTypeHillsUphillFinish"
	case pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_FLAT_FINISH:
		messageKey = "StageTypeMountainsFlatFinish"
	case pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_UPHILL_FINISH:
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
	var messages []string
	messages = append(messages, localizer.localize(localizeParams{
		key: "GapMinutes",
		data: map[string]interface{}{
			"Minutes": minutes,
		},
		pluralCount: minutes,
	}))
	if seconds > 0 {
		messages = append(messages, localizer.localize(localizeParams{
			key: "GapAndSeconds",
			data: map[string]interface{}{
				"Seconds": seconds,
			},
			pluralCount: seconds,
		}))
	}
	return strings.Join(messages, " ")
}
