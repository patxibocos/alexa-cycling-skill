package cycling

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"time"
)

type RaceResult interface {
	isRaceResult()
}

type RaceStage interface {
	isRaceStage()
}

type MountainsStage interface {
	isMountainsStage()
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
type MultiTTTStageRaceWithResults struct {
	StageNumber int
	Top3Teams   *Top3Teams
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

type Top3Teams struct {
	First  *TeamResult
	Second *TeamResult
	Third  *TeamResult
}

type TeamResult struct {
	Team *pcsscraper.Team
	Time int64
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
func (_ MultiTTTStageRaceWithResults) isRaceResult() {}
func (_ MultiStageRaceWithoutResults) isRaceResult() {}

type NoStage struct{}
type StageWithData struct {
	Departure      string
	Arrival        string
	Distance       float32
	ProfileType    pcsscraper.Stage_ProfileType
	StartDateLocal time.Time
	StageType      pcsscraper.Stage_StageType
}
type StageWithoutData struct {
	StartDateLocal time.Time
}

func (_ RestDayStage) isRaceStage()     {}
func (_ NoStage) isRaceStage()          {}
func (_ StageWithData) isRaceStage()    {}
func (_ StageWithoutData) isRaceStage() {}

type SingleDayRace struct{}
type YesMountainsStage struct {
	StageNumber    int
	StartDateLocal time.Time
}
type NoStageTypeData struct{}
type NoMountainsStage struct{}

func (_ SingleDayRace) isMountainsStage()     {}
func (_ YesMountainsStage) isMountainsStage() {}
func (_ NoStageTypeData) isMountainsStage()   {}
func (_ NoMountainsStage) isMountainsStage()  {}
