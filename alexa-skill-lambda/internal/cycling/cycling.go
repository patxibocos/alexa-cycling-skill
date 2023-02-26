package cycling

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/timeutils"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"time"
)

const MillisModulo = 60000

func GetRaceResult(race *pcsscraper.Race, riders []*pcsscraper.Rider, teams []*pcsscraper.Team, location *time.Location) RaceResult {
	if raceIsFromThePast(race, location) {
		return &PastRace{GcTop3: GetTop3FromResult(race.Result, riders)}
	}
	if raceIsFromTheFuture(race, location) {
		return new(FutureRace)
	}
	if IsSingleDayRace(race, location) {
		if areStageResultsAvailable(race.Stages[0]) {
			return &SingleDayRaceWithResults{
				Top3: GetTop3FromResult(race.Stages[0].Result, riders),
			}
		} else {
			return new(SingleDayRaceWithoutResults)
		}
	}
	todayStage, stageNumber := findTodayStage(race, location)
	if todayStage == nil {
		return new(RestDayStage)
	}
	if areStageResultsAvailable(todayStage) {
		if isTeamTimeTrial(todayStage) {
			return &MultiTTTStageRaceWithResults{
				Top3Teams:   GetTop3TeamsFromResult(todayStage.Result, teams),
				GcTop3:      GetTop3FromResult(race.Result, riders),
				StageNumber: stageNumber,
			}
		}
		return &MultiStageRaceWithResults{
			Top3:        GetTop3FromResult(todayStage.Result, riders),
			GcTop3:      GetTop3FromResult(race.Result, riders),
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

func GetActiveRaces(races []*pcsscraper.Race, location *time.Location) []*pcsscraper.Race {
	var activeRaces []*pcsscraper.Race
	for _, race := range races {
		if raceIsActive(race, location) {
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

func FindNextRace(races []*pcsscraper.Race, location *time.Location) *pcsscraper.Race {
	today := timeutils.Today(location)
	for _, race := range races {
		if race.StartDateLocal(location).After(today) {
			return race
		}
	}
	return nil
}

func FindMountainsStage(race *pcsscraper.Race, location *time.Location) MountainsStage {
	if IsSingleDayRace(race, location) {
		return new(SingleDayRace)
	}
	for i, stage := range race.Stages {
		if stage.ProfileType == pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_FLAT_FINISH || stage.ProfileType == pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_UPHILL_FINISH {
			return &YesMountainsStage{
				StageNumber:    i + 1,
				StartDateLocal: stage.StartDateTimeLocal(location),
			}
		}
	}
	if race.Stages[0].ProfileType == pcsscraper.Stage_PROFILE_TYPE_UNSPECIFIED {
		return new(NoStageTypeData)
	}
	return new(NoMountainsStage)
}

func GetRaceStageForDay(race *pcsscraper.Race, day time.Time, location *time.Location) RaceStage {
	var stage *pcsscraper.Stage
	dayYear, dayMonth, dayDay := day.Date()
	for _, s := range race.Stages {
		stageYear, stageMonth, stageDay := s.StartDateTimeLocal(location).Date()
		if stageYear == dayYear && stageMonth == dayMonth && stageDay == dayDay {
			stage = s
		}
	}
	if stage == nil {
		// Check if rest day
		if race.StartDateLocal(location).Before(day) && race.EndDateLocal(location).After(day) {
			return new(RestDayStage)
		}
		return new(NoStage)
	}
	if !StageContainsData(stage) {
		return &StageWithoutData{
			StartDateLocal: stage.StartDateTimeLocal(location),
		}
	}
	return &StageWithData{
		Departure:      stage.GetDeparture(),
		Arrival:        stage.GetArrival(),
		Distance:       stage.GetDistance(),
		ProfileType:    stage.GetProfileType(),
		StartDateLocal: stage.StartDateTimeLocal(location),
		StageType:      stage.StageType,
	}
}

func GetRaceStageForIndex(race *pcsscraper.Race, index int, location *time.Location) RaceStage {
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
			StartDateLocal: stage.StartDateTimeLocal(location),
		}
	}
	return &StageWithData{
		Departure:      stage.GetDeparture(),
		Arrival:        stage.GetArrival(),
		Distance:       stage.GetDistance(),
		ProfileType:    stage.GetProfileType(),
		StartDateLocal: stage.StartDateTimeLocal(location),
		StageType:      stage.GetStageType(),
	}
}

func StageContainsData(stage *pcsscraper.Stage) bool {
	return (stage.GetDeparture() != "" && stage.GetArrival() != "") || stage.GetDistance() > 0 || stage.GetProfileType() != pcsscraper.Stage_PROFILE_TYPE_UNSPECIFIED
}

func findRider(riderID string, riders []*pcsscraper.Rider) *pcsscraper.Rider {
	for _, rider := range riders {
		if riderID == rider.Id {
			return rider
		}
	}
	return nil
}

func findTeam(teamID string, teams []*pcsscraper.Team) *pcsscraper.Team {
	for _, team := range teams {
		if teamID == team.Id {
			return team
		}
	}
	return nil
}

func GetTop3FromResult(result []*pcsscraper.ParticipantResult, riders []*pcsscraper.Rider) *Top3 {
	return &Top3{
		First: &RiderResult{
			Rider: findRider(result[0].ParticipantId, riders),
			Time:  result[0].Time,
		},
		Second: &RiderResult{
			Rider: findRider(result[1].ParticipantId, riders),
			Time:  result[1].Time,
		},
		Third: &RiderResult{
			Rider: findRider(result[2].ParticipantId, riders),
			Time:  result[2].Time,
		},
	}
}

func GetTop3TeamsFromResult(result []*pcsscraper.ParticipantResult, teams []*pcsscraper.Team) *Top3Teams {
	return &Top3Teams{
		First: &TeamResult{
			Team: findTeam(result[0].ParticipantId, teams),
			Time: result[0].Time,
		},
		Second: &TeamResult{
			Team: findTeam(result[1].ParticipantId, teams),
			Time: result[1].Time,
		},
		Third: &TeamResult{
			Team: findTeam(result[2].ParticipantId, teams),
			Time: result[2].Time,
		},
	}
}

func raceIsActive(race *pcsscraper.Race, location *time.Location) bool {
	today := timeutils.Today(location)
	return (race.StartDateLocal(location) == today || race.StartDateLocal(location).Before(today)) &&
		(race.EndDateLocal(location) == today || race.EndDateLocal(location).After(today))
}

func findTodayStage(race *pcsscraper.Race, location *time.Location) (*pcsscraper.Stage, int) {
	today := timeutils.Today(location)
	for i, stage := range race.Stages {
		stageYear, stageMonth, stageDay := stage.StartDateTimeLocal(location).Date()
		todayYear, todayMonth, todayDay := today.Date()
		if stageYear == todayYear && stageMonth == todayMonth && stageDay == todayDay {
			return stage, i + 1
		}
	}
	return nil, 0
}

func raceIsFromThePast(race *pcsscraper.Race, location *time.Location) bool {
	return race.EndDateLocal(location).Before(timeutils.Today(location))
}

func raceIsFromTheFuture(race *pcsscraper.Race, location *time.Location) bool {
	return race.StartDateLocal(location).After(timeutils.Today(location))
}

func RaceIsFinished(race *pcsscraper.Race) bool {
	return areStageResultsAvailable(race.Stages[len(race.Stages)-1])
}

func RaceHasNotStarted(race *pcsscraper.Race) bool {
	return !areStageResultsAvailable(race.Stages[0])
}

func IsSingleDayRace(race *pcsscraper.Race, location *time.Location) bool {
	return race.StartDateLocal(location) == race.EndDateLocal(location)
}

func areStageResultsAvailable(stage *pcsscraper.Stage) bool {
	return stage.Result != nil && len(stage.Result) >= 3
}

func isTeamTimeTrial(stage *pcsscraper.Stage) bool {
	return stage.StageType == pcsscraper.Stage_STAGE_TYPE_TEAM_TIME_TRIAL
}

func MillisForRace(race *pcsscraper.Race) int {
	raceIdSum := sha1.Sum([]byte(race.Id))
	sumBytes := make([]byte, len(raceIdSum))
	copy(sumBytes, raceIdSum[:])
	millisFromRaceId := int(binary.BigEndian.Uint16(sumBytes) % MillisModulo)
	return millisFromRaceId
}
