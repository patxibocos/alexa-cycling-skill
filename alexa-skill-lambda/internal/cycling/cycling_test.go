package cycling

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

const Day = 24 * time.Hour

func TestPastRace(t *testing.T) {
	yesterday := today().Add(-1 * Day)
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(yesterday),
		EndDate:   timestamppb.New(yesterday),
		Result:    []*pcsscraper.RiderResult{{RiderId: "ID1"}, {RiderId: "ID2"}, {RiderId: "ID3"}},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders)

	assert.Equal(t, &PastRace{
		GcTop3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
	}, raceResult)
}

func TestFutureRace(t *testing.T) {
	tomorrow := today().Add(1 * Day)
	dayAfterTomorrow := today().Add(2 * Day)
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(tomorrow),
		EndDate:   timestamppb.New(dayAfterTomorrow),
	}

	raceResult := GetRaceResult(race, nil)

	assert.Equal(t, new(FutureRace), raceResult)
}

func TestMultiStageRaceWithoutResults(t *testing.T) {
	today := today()
	yesterday := today.Add(-1 * Day)
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(yesterday),
		EndDate:   timestamppb.New(today),
		Stages: []*pcsscraper.Stage{
			{StartDate: timestamppb.New(yesterday)},
			{StartDate: timestamppb.New(today)},
		},
	}

	raceResult := GetRaceResult(race, nil)

	assert.Equal(t, &MultiStageRaceWithoutResults{
		StageNumber: 2,
	}, raceResult)
}

func TestMultiStageRaceWithResults(t *testing.T) {
	today := today()
	yesterday := today.Add(-1 * Day)
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(yesterday),
		EndDate:   timestamppb.New(today),
		Stages: []*pcsscraper.Stage{
			{StartDate: timestamppb.New(yesterday)},
			{StartDate: timestamppb.New(today), Result: []*pcsscraper.RiderResult{{RiderId: "ID1"}, {RiderId: "ID2"}, {RiderId: "ID3"}}},
		},
		Result: []*pcsscraper.RiderResult{{RiderId: "ID1"}, {RiderId: "ID2"}, {RiderId: "ID3"}},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders)

	assert.Equal(t, &MultiStageRaceWithResults{
		Top3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
		GcTop3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
		StageNumber: 2,
		IsLastStage: true,
	}, raceResult)
}

func TestSingleDayRaceWithoutResults(t *testing.T) {
	today := today()
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(today),
		EndDate:   timestamppb.New(today),
		Stages: []*pcsscraper.Stage{
			{StartDate: timestamppb.New(today)},
		},
	}

	raceResult := GetRaceResult(race, nil)

	assert.Equal(t, new(SingleDayRaceWithoutResults), raceResult)
}

func TestSingleDayRaceWithResults(t *testing.T) {
	today := today()
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(today),
		EndDate:   timestamppb.New(today),
		Stages: []*pcsscraper.Stage{
			{StartDate: timestamppb.New(today), Result: []*pcsscraper.RiderResult{{RiderId: "ID1"}, {RiderId: "ID2"}, {RiderId: "ID3"}}},
		},
		Result: []*pcsscraper.RiderResult{{RiderId: "ID1"}, {RiderId: "ID2"}, {RiderId: "ID3"}},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders)

	assert.Equal(t, &SingleDayRaceWithResults{
		Top3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
	}, raceResult)
}

func TestRestDayStage(t *testing.T) {
	yesterday := today().Add(-1 * Day)
	tomorrow := today().Add(1 * Day)
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(yesterday),
		EndDate:   timestamppb.New(tomorrow),
		Stages:    []*pcsscraper.Stage{{StartDate: timestamppb.New(yesterday)}, {StartDate: timestamppb.New(tomorrow)}},
	}

	raceResult := GetRaceResult(race, nil)

	assert.Equal(t, new(RestDayStage), raceResult)
}

func TestFindNextRaceNotFound(t *testing.T) {
	yesterday := timestamppb.New(today().Add(-1 * Day))
	today := timestamppb.New(today())
	races := []*pcsscraper.Race{
		{StartDate: yesterday},
		{StartDate: today},
	}

	nextRace := FindNextRace(races)

	assert.Nil(t, nextRace)
}

func TestFindNextRaceIsFound(t *testing.T) {
	tomorrow := timestamppb.New(today().Add(1 * Day))
	races := []*pcsscraper.Race{
		{StartDate: tomorrow},
	}

	nextRace := FindNextRace(races)

	assert.NotNil(t, nextRace)
}

func TestGetActiveRaces(t *testing.T) {
	yesterday := today().Add(-1 * Day)
	tomorrow := today().Add(1 * Day)
	today := today()
	races := []*pcsscraper.Race{
		{StartDate: timestamppb.New(yesterday), EndDate: timestamppb.New(today)},
		{StartDate: timestamppb.New(today), EndDate: timestamppb.New(today)},
		{StartDate: timestamppb.New(tomorrow), EndDate: timestamppb.New(tomorrow)},
	}

	activeRaces := GetActiveRaces(races)

	assert.Len(t, activeRaces, 2)
}

func TestNoStage(t *testing.T) {
	tomorrow := today().Add(1 * Day)
	today := today()
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{StartDate: timestamppb.New(tomorrow)}}}

	raceStage := GetRaceStage(race, today)

	assert.Equal(t, new(NoStage), raceStage)
}

func TestRestDay(t *testing.T) {
	yesterday := today().Add(-1 * Day)
	tomorrow := today().Add(1 * Day)
	today := today()
	race := &pcsscraper.Race{
		StartDate: timestamppb.New(yesterday),
		EndDate:   timestamppb.New(tomorrow),
		Stages:    []*pcsscraper.Stage{{StartDate: timestamppb.New(yesterday)}, {StartDate: timestamppb.New(tomorrow)}},
	}

	raceStage := GetRaceStage(race, today)

	assert.Equal(t, new(RestDayStage), raceStage)
}

func TestStageWithoutData(t *testing.T) {
	today := today()
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{StartDate: timestamppb.New(today)}}}

	raceStage := GetRaceStage(race, today)

	assert.Equal(t, new(StageWithoutData), raceStage)
}

func TestStageWithData(t *testing.T) {
	today := today()
	bilbao := "Bilbao"
	barcelona := "Barcelona"
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{
		StartDate: timestamppb.New(today),
		Departure: &bilbao,
		Arrival:   &barcelona,
		Distance:  123.456,
		Type:      pcsscraper.Stage_TYPE_MOUNTAINS_UPHILL_FINISH,
	}}}

	raceStage := GetRaceStage(race, today)

	assert.Equal(t, &StageWithData{
		Departure: "Bilbao",
		Arrival:   "Barcelona",
		Distance:  123.456,
		Type:      pcsscraper.Stage_TYPE_MOUNTAINS_UPHILL_FINISH,
	}, raceStage)
}
