package cycling

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"time"
)

const MillisModulo = 60000

func GetRaceResult(race *pcsscraper.Race, riders []*pcsscraper.Rider) RaceResult {
	if raceIsFromThePast(race) {
		return &PastRace{GcTop3: getTop3FromResult(race.Result, riders)}
	}
	if raceIsFromTheFuture(race) {
		return new(FutureRace)
	}
	if IsSingleDayRace(race) {
		if areStageResultsAvailable(race.Stages[0]) {
			return &SingleDayRaceWithResults{
				Top3: getTop3FromResult(race.Stages[0].Result, riders),
			}
		} else {
			return new(SingleDayRaceWithoutResults)
		}
	}
	todayStage, stageNumber := findTodayStage(race)
	if todayStage == nil {
		return new(RestDayStage)
	}
	if areStageResultsAvailable(todayStage) {
		return &MultiStageRaceWithResults{
			Top3:        getTop3FromResult(todayStage.Result, riders),
			GcTop3:      getTop3FromResult(race.Result, riders),
			StageNumber: stageNumber,
		}
	}
	return &MultiStageRaceWithoutResults{
		StageNumber: stageNumber,
	}
}

func IsLastRaceStage(race *pcsscraper.Race, stageNumber int) bool {
	return stageNumber == len(race.Stages)
}

func GetActiveRaces(races []*pcsscraper.Race) []*pcsscraper.Race {
	var activeRaces []*pcsscraper.Race
	for _, race := range races {
		if raceIsActive(race) {
			activeRaces = append(activeRaces, race)
		}
	}
	return activeRaces
}

func FindRace(races []*pcsscraper.Race, raceId string) *pcsscraper.Race {
	var race *pcsscraper.Race
	for _, r := range races {
		if r.Id == raceId {
			race = r
		}
	}
	return race
}

func FindNextRace(races []*pcsscraper.Race) *pcsscraper.Race {
	today := today()
	for _, race := range races {
		if race.StartDate.AsTime().After(today) {
			return race
		}
	}
	return nil
}

func FindMountainsStage(race *pcsscraper.Race) MountainsStage {
	if IsSingleDayRace(race) {
		return new(SingleDayRace)
	}
	for i, stage := range race.Stages {
		if stage.Type == pcsscraper.Stage_TYPE_MOUNTAINS_FLAT_FINISH || stage.Type == pcsscraper.Stage_TYPE_MOUNTAINS_UPHILL_FINISH {
			return &YesMountainsStage{
				StageNumber: i + 1,
				StartDate:   stage.StartDateTime,
			}
		}
	}
	if race.Stages[0].Type == pcsscraper.Stage_TYPE_UNSPECIFIED {
		return new(NoStageTypeData)
	}
	return new(NoMountainsStage)
}

func GetRaceStageForDay(race *pcsscraper.Race, day time.Time) RaceStage {
	var stage *pcsscraper.Stage
	for _, s := range race.Stages {
		if s.StartDateTime.AsTime() == day {
			stage = s
		}
	}
	if stage == nil {
		// Check if rest day
		if race.StartDate.AsTime().Before(day) && race.EndDate.AsTime().After(day) {
			return new(RestDayStage)
		}
		return new(NoStage)
	}
	if !StageContainsData(stage) {
		return &StageWithoutData{
			StartDate: stage.StartDateTime,
		}
	}
	return &StageWithData{
		Departure: stage.GetDeparture(),
		Arrival:   stage.GetArrival(),
		Distance:  stage.GetDistance(),
		Type:      stage.GetType(),
		StartDate: stage.StartDateTime,
	}
}

func GetRaceStageForIndex(race *pcsscraper.Race, index int) RaceStage {
	var stage *pcsscraper.Stage
	for i, s := range race.Stages {
		if i+1 == index {
			stage = s
		}
	}
	if stage == nil {
		return new(NoStage)
	}
	if !StageContainsData(stage) {
		return &StageWithoutData{
			StartDate: stage.GetStartDateTime(),
		}
	}
	return &StageWithData{
		Departure: stage.GetDeparture(),
		Arrival:   stage.GetArrival(),
		Distance:  stage.GetDistance(),
		Type:      stage.GetType(),
		StartDate: stage.GetStartDateTime(),
	}
}

func StageContainsData(stage *pcsscraper.Stage) bool {
	return (stage.GetDeparture() != "" && stage.GetArrival() != "") || stage.GetDistance() > 0 || stage.GetType() != pcsscraper.Stage_TYPE_UNSPECIFIED
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
		First: &RiderResult{
			Rider: findRider(result[0].RiderId, riders),
			Time:  result[0].Time,
		},
		Second: &RiderResult{
			Rider: findRider(result[1].RiderId, riders),
			Time:  result[1].Time,
		},
		Third: &RiderResult{
			Rider: findRider(result[2].RiderId, riders),
			Time:  result[2].Time,
		},
	}
}

func raceIsActive(race *pcsscraper.Race) bool {
	today := today()
	return (race.StartDate.AsTime() == today || race.StartDate.AsTime().Before(today)) &&
		(race.EndDate.AsTime() == today || race.EndDate.AsTime().After(today))
}

func findTodayStage(race *pcsscraper.Race) (*pcsscraper.Stage, int) {
	today := today()
	for i, stage := range race.Stages {
		if stage.StartDateTime.AsTime() == today {
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

func IsSingleDayRace(race *pcsscraper.Race) bool {
	return race.StartDate.AsTime() == race.EndDate.AsTime()
}

func areStageResultsAvailable(stage *pcsscraper.Stage) bool {
	return stage.Result != nil && len(stage.Result) > 0
}

func today() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today
}

func MillisForRace(race *pcsscraper.Race) int {
	raceIdSum := sha1.Sum([]byte(race.Id))
	sumBytes := make([]byte, len(raceIdSum))
	copy(sumBytes, raceIdSum[:])
	millisFromRaceId := int(binary.BigEndian.Uint16(sumBytes) % MillisModulo)
	return millisFromRaceId
}
