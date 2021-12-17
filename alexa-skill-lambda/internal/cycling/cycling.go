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

func findTodayStage(race *pcsscraper.Race) (*pcsscraper.Stage, int) {
	today := today()
	for i, stage := range race.Stages {
		if stage.StartDate.AsTime() == today {
			return stage, i + 1
		}
	}
	return nil, 0
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
	todayStage, stageNumber := findTodayStage(race)
	if todayStage == nil {
		return buildRestDayStage()
	}
	if areStageResultsAvailable(todayStage) {
		return buildMultiStageRaceWithResults(race, todayStage, stageNumber, riders)
	}
	return buildMultiStageRaceWithoutResults(stageNumber)
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

func findRider(riderID string, riders []*pcsscraper.Rider) *pcsscraper.Rider {
	for _, rider := range riders {
		if riderID == rider.Id {
			return rider
		}
	}
	return nil
}

func getTop3FromResult(result []*pcsscraper.RiderResult, riders []*pcsscraper.Rider) *Top3 {
	return &Top3{
		First:  findRider(result[0].RiderId, riders),
		Second: findRider(result[1].RiderId, riders),
		Third:  findRider(result[2].RiderId, riders),
	}
}

func buildFutureRace() *FutureRace {
	return new(FutureRace)
}

func buildSingleDayRaceWithResults(race *pcsscraper.Race, riders []*pcsscraper.Rider) *SingleDayRaceWithResults {
	return &SingleDayRaceWithResults{
		Top3: getTop3FromResult(race.Stages[0].Result, riders),
	}
}

func buildMultiStageRaceWithResults(race *pcsscraper.Race, stage *pcsscraper.Stage, stageNumber int, riders []*pcsscraper.Rider) *MultiStageRaceWithResults {
	return &MultiStageRaceWithResults{
		Top3:        getTop3FromResult(stage.Result, riders),
		GcTop3:      getTop3FromResult(race.Result, riders),
		StageNumber: stageNumber,
		IsLastStage: stageNumber == len(race.Stages),
	}
}

func buildMultiStageRaceWithoutResults(stageNumber int) *MultiStageRaceWithoutResults {
	return &MultiStageRaceWithoutResults{
		StageNumber: stageNumber,
	}
}

type RaceInfo interface {
	isRaceInfo()
}

type PastRace struct{ GcTop3 *Top3 }
type FutureRace struct{}
type RestDayStage struct{}
type SingleDayRaceWithResults struct{ Top3 *Top3 }
type SingleDayRaceWithoutResults struct{}
type MultiStageRaceWithResults struct {
	StageNumber int
	Top3        *Top3
	GcTop3      *Top3
	IsLastStage bool
}
type MultiStageRaceWithoutResults struct {
	StageNumber int
}

type Top3 struct {
	First  *pcsscraper.Rider
	Second *pcsscraper.Rider
	Third  *pcsscraper.Rider
}

func (_ PastRace) isRaceInfo()                     {}
func (_ FutureRace) isRaceInfo()                   {}
func (_ RestDayStage) isRaceInfo()                 {}
func (_ SingleDayRaceWithResults) isRaceInfo()     {}
func (_ SingleDayRaceWithoutResults) isRaceInfo()  {}
func (_ MultiStageRaceWithResults) isRaceInfo()    {}
func (_ MultiStageRaceWithoutResults) isRaceInfo() {}
