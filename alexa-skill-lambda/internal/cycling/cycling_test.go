package cycling

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/timeutils"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

const Day = 24 * time.Hour

func TestPastRace(t *testing.T) {
	yesterday := timeutils.Today(time.UTC).Add(-1 * Day)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{{
			StartDateTime: timestamppb.New(yesterday),
			GeneralResults: &pcsscraper.GeneralResults{
				Time: []*pcsscraper.ParticipantResultTime{{ParticipantId: "ID1"}, {ParticipantId: "ID2"}, {ParticipantId: "ID3"}},
			},
		}},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders, nil, time.UTC)

	assert.Equal(t, &PastRace{
		GcTop3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
	}, raceResult)
}

func TestFutureRace(t *testing.T) {
	tomorrow := timeutils.Today(time.UTC).Add(1 * Day)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{{
			StartDateTime: timestamppb.New(tomorrow),
		}},
	}

	raceResult := GetRaceResult(race, nil, nil, time.UTC)

	assert.Equal(t, new(FutureRace), raceResult)
}

func TestMultiStageRaceWithoutResults(t *testing.T) {
	today := timeutils.Today(time.UTC)
	yesterday := today.Add(-1 * Day)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{
			{StartDateTime: timestamppb.New(yesterday)},
			{StartDateTime: timestamppb.New(today)},
		},
	}

	raceResult := GetRaceResult(race, nil, nil, time.UTC)

	assert.Equal(t, &MultiStageRaceWithoutResults{
		StageNumber: 2,
	}, raceResult)
}

func TestMultiStageRaceWithResults(t *testing.T) {
	today := timeutils.Today(time.UTC)
	yesterday := today.Add(-1 * Day)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{
			{StartDateTime: timestamppb.New(yesterday)},
			{
				StartDateTime: timestamppb.New(today),
				StageResults: &pcsscraper.StageResults{
					Time: []*pcsscraper.ParticipantResultTime{{ParticipantId: "ID1"}, {ParticipantId: "ID2"}, {ParticipantId: "ID3"}},
				},
				GeneralResults: &pcsscraper.GeneralResults{
					Time: []*pcsscraper.ParticipantResultTime{{ParticipantId: "ID1"}, {ParticipantId: "ID2"}, {ParticipantId: "ID3"}},
				},
			},
		},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders, nil, time.UTC)

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
	}, raceResult)
}

func TestSingleDayRaceWithoutResults(t *testing.T) {
	today := timeutils.Today(time.UTC)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{
			{StartDateTime: timestamppb.New(today)},
		},
	}

	raceResult := GetRaceResult(race, nil, nil, time.UTC)

	assert.Equal(t, new(SingleDayRaceWithoutResults), raceResult)
}

func TestSingleDayRaceWithResults(t *testing.T) {
	today := timeutils.Today(time.UTC)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{
			{
				StartDateTime: timestamppb.New(today),
				StageResults: &pcsscraper.StageResults{
					Time: []*pcsscraper.ParticipantResultTime{{ParticipantId: "ID1"}, {ParticipantId: "ID2"}, {ParticipantId: "ID3"}},
				},
				GeneralResults: &pcsscraper.GeneralResults{
					Time: []*pcsscraper.ParticipantResultTime{{ParticipantId: "ID1"}, {ParticipantId: "ID2"}, {ParticipantId: "ID3"}},
				},
			},
		},
	}
	riders := []*pcsscraper.Rider{{Id: "ID1"}, {Id: "ID2"}, {Id: "ID3"}}

	raceResult := GetRaceResult(race, riders, nil, time.UTC)

	assert.Equal(t, &SingleDayRaceWithResults{
		Top3: &Top3{
			First:  &RiderResult{Rider: riders[0]},
			Second: &RiderResult{Rider: riders[1]},
			Third:  &RiderResult{Rider: riders[2]},
		},
	}, raceResult)
}

func TestRestDayStage(t *testing.T) {
	yesterday := timeutils.Today(time.UTC).Add(-1 * Day)
	tomorrow := timeutils.Today(time.UTC).Add(1 * Day)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(yesterday)}, {StartDateTime: timestamppb.New(tomorrow)}},
	}

	raceResult := GetRaceResult(race, nil, nil, time.UTC)

	assert.Equal(t, new(RestDayStage), raceResult)
}

func TestFindNextRaceNotFound(t *testing.T) {
	yesterday := timestamppb.New(timeutils.Today(time.UTC).Add(-1 * Day))
	today := timestamppb.New(timeutils.Today(time.UTC))
	races := []*pcsscraper.Race{
		{Stages: []*pcsscraper.Stage{{StartDateTime: yesterday}}},
		{Stages: []*pcsscraper.Stage{{StartDateTime: today}}},
	}

	nextRace := FindNextRace(races, time.UTC)

	assert.Nil(t, nextRace)
}

func TestFindNextRaceIsFound(t *testing.T) {
	tomorrow := timestamppb.New(timeutils.Today(time.UTC).Add(1 * Day))
	races := []*pcsscraper.Race{
		{Stages: []*pcsscraper.Stage{{StartDateTime: tomorrow}}},
	}

	nextRace := FindNextRace(races, time.UTC)

	assert.NotNil(t, nextRace)
}

func TestGetActiveRaces(t *testing.T) {
	yesterday := timeutils.Today(time.UTC).Add(-1 * Day)
	tomorrow := timeutils.Today(time.UTC).Add(1 * Day)
	today := timeutils.Today(time.UTC)
	races := []*pcsscraper.Race{
		{Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(yesterday)}, {StartDateTime: timestamppb.New(today)}}},
		{Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(today)}, {StartDateTime: timestamppb.New(today)}}},
		{Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(tomorrow)}, {StartDateTime: timestamppb.New(tomorrow)}}},
	}

	activeRaces := GetActiveRaces(races, time.UTC)

	assert.Len(t, activeRaces, 2)
}

func TestNoStage(t *testing.T) {
	tomorrow := timeutils.Today(time.UTC).Add(1 * Day)
	today := timeutils.Today(time.UTC)
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(tomorrow)}}}

	raceStage := GetRaceStageForDay(race, today, time.UTC)

	assert.Equal(t, new(NoStage), raceStage)
}

func TestRestDay(t *testing.T) {
	yesterday := timeutils.Today(time.UTC).Add(-1 * Day)
	tomorrow := timeutils.Today(time.UTC).Add(1 * Day)
	today := timeutils.Today(time.UTC)
	race := &pcsscraper.Race{
		Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(yesterday)}, {StartDateTime: timestamppb.New(tomorrow)}},
	}

	raceStage := GetRaceStageForDay(race, today, time.UTC)

	assert.Equal(t, new(RestDayStage), raceStage)
}

func TestStageWithoutData(t *testing.T) {
	today := timeutils.Today(time.UTC)
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{StartDateTime: timestamppb.New(today)}}}

	raceStage := GetRaceStageForDay(race, today, time.UTC)

	assert.Equal(t, &StageWithoutData{StartDateLocal: today}, raceStage)
}

func TestStageWithData(t *testing.T) {
	today := timeutils.Today(time.UTC)
	bilbao := "Bilbao"
	barcelona := "Barcelona"
	race := &pcsscraper.Race{Stages: []*pcsscraper.Stage{{
		StartDateTime: timestamppb.New(today),
		Departure:     &bilbao,
		Arrival:       &barcelona,
		Distance:      123.456,
		ProfileType:   pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_UPHILL_FINISH,
	}}}

	raceStage := GetRaceStageForDay(race, today, time.UTC)

	assert.Equal(t, &StageWithData{
		Departure:      "Bilbao",
		Arrival:        "Barcelona",
		Distance:       123.456,
		ProfileType:    pcsscraper.Stage_PROFILE_TYPE_MOUNTAINS_UPHILL_FINISH,
		StartDateLocal: today,
	}, raceStage)
}
