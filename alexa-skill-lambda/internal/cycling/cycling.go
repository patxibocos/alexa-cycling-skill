package cycling

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"time"
)

func GetNextRace(races []*pcsscraper.Race) *pcsscraper.Race {
	now := time.Now()
	for _, race := range races {
		if race.StartDate.AsTime().After(now) {
			return race
		}
	}
	return nil
}

func today() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today
}

type RaceResult struct {
	top3   []*pcsscraper.Rider
	gcTop3 []*pcsscraper.Rider
}

func findTodayStage(race *pcsscraper.Race) *pcsscraper.Stage {
	today := today()
	for _, stage := range race.Stages {
		if stage.StartDate.AsTime() == today {
			return stage
		}
	}
	return nil
}

func raceIsFromThePast(race *pcsscraper.Race) bool {
	return race.EndDate.AsTime().Before(today())
}

func raceIsFromTheFuture(race *pcsscraper.Race) bool {
	return race.StartDate.AsTime().After(today())
}

func isSingleDayRace(race *pcsscraper.Race) bool {
	return race.StartDate.AsTime() == race.EndDate.AsTime()
}

func areStageResultsAvailable(stage *pcsscraper.Stage) bool {
	return stage.Result != nil && len(stage.Result) > 0
}

func GetRaceResult(race *pcsscraper.Race, riders []*pcsscraper.Rider) RaceInfo {
	if raceIsFromThePast(race) {
		return buildPastRace(race, riders)
	}
	if raceIsFromTheFuture(race) {
		return buildFutureRace()
	}
	if isSingleDayRace(race) {
		if areStageResultsAvailable(race.Stages[0]) {
			return buildSingleDayRaceWithResults(race, riders)
		} else {
			return buildSingleDayRaceWithoutResults()
		}
	}
	todayStage := findTodayStage(race)
	if todayStage == nil {
		return buildRestDayStage()
	}
	if areStageResultsAvailable(todayStage) {
		return buildMultiStageRaceWithResults(race, todayStage, riders)
	}
	return buildMultiStageRaceWithoutResults()
}

func buildSingleDayRaceWithoutResults() *SingleDayRaceWithoutResults {
	return new(SingleDayRaceWithoutResults)
}

func buildRestDayStage() *RestDayStage {
	return new(RestDayStage)
}

func buildPastRace(race *pcsscraper.Race, riders []*pcsscraper.Rider) *PastRace {
	return &PastRace{
		GcTop3: getTop3FromResult(race.Result, riders),
	}
}

func getTop3FromResult(result []*pcsscraper.RiderResult, riders []*pcsscraper.Rider) []*pcsscraper.Rider {
	var top3 []*pcsscraper.Rider
	riderIDs := []string{result[0].RiderId, result[1].RiderId, result[2].RiderId}
	for _, rider := range riders {
		for _, riderID := range riderIDs {
			if riderID == rider.Id {
				top3 = append(top3, rider)
			}
		}
		if len(top3) == len(riderIDs) {
			break
		}
	}
	return top3
}

func buildFutureRace() *FutureRace {
	return new(FutureRace)
}

func buildSingleDayRaceWithResults(race *pcsscraper.Race, riders []*pcsscraper.Rider) *SingleDayRaceWithResults {
	return &SingleDayRaceWithResults{
		Top3: getTop3FromResult(race.Stages[0].Result, riders),
	}
}

func buildMultiStageRaceWithResults(race *pcsscraper.Race, stage *pcsscraper.Stage, riders []*pcsscraper.Rider) *MultiStageRaceWithResults {
	return &MultiStageRaceWithResults{
		Top3:   getTop3FromResult(stage.Result, riders),
		GcTop3: getTop3FromResult(race.Result, riders),
	}
}

func buildMultiStageRaceWithoutResults() *MultiStageRaceWithoutResults {
	return new(MultiStageRaceWithoutResults)
}

type RaceInfo interface {
	isRaceInfo()
}

type PastRace struct{ GcTop3 []*pcsscraper.Rider }
type FutureRace struct{}
type RestDayStage struct{}
type SingleDayRaceWithResults struct{ Top3 []*pcsscraper.Rider }
type SingleDayRaceWithoutResults struct{}
type MultiStageRaceWithResults struct {
	Top3   []*pcsscraper.Rider
	GcTop3 []*pcsscraper.Rider
}
type MultiStageRaceWithoutResults struct{}

func (_ PastRace) isRaceInfo()                     {}
func (_ FutureRace) isRaceInfo()                   {}
func (_ RestDayStage) isRaceInfo()                 {}
func (_ SingleDayRaceWithResults) isRaceInfo()     {}
func (_ SingleDayRaceWithoutResults) isRaceInfo()  {}
func (_ MultiStageRaceWithResults) isRaceInfo()    {}
func (_ MultiStageRaceWithoutResults) isRaceInfo() {}
