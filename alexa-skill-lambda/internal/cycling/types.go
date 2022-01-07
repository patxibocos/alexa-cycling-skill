package cycling

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RaceResult interface {
	isRaceResult()
}

type RaceStage interface {
	isRaceStage()
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
}
type MultiStageRaceWithoutResults struct {
	StageNumber int
}

type Top3 struct {
	First  *RiderResult
	Second *RiderResult
	Third  *RiderResult
}

type RiderResult struct {
	Rider *pcsscraper.Rider
	Time  int64
}

func (_ PastRace) isRaceResult()                     {}
func (_ FutureRace) isRaceResult()                   {}
func (_ RestDayStage) isRaceResult()                 {}
func (_ SingleDayRaceWithResults) isRaceResult()     {}
func (_ SingleDayRaceWithoutResults) isRaceResult()  {}
func (_ MultiStageRaceWithResults) isRaceResult()    {}
func (_ MultiStageRaceWithoutResults) isRaceResult() {}

type NoStage struct{}
type StageWithData struct {
	Departure string
	Arrival   string
	Distance  float32
	Type      pcsscraper.Stage_Type
	StartDate *timestamppb.Timestamp
}
type StageWithoutData struct {
	StartDate *timestamppb.Timestamp
}

func (_ RestDayStage) isRaceStage()     {}
func (_ NoStage) isRaceStage()          {}
func (_ StageWithData) isRaceStage()    {}
func (_ StageWithoutData) isRaceStage() {}
