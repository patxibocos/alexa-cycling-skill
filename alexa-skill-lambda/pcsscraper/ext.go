package pcsscraper

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/timeutils"
	"time"
)

func (r *Race) StartDateLocal(location *time.Location) time.Time {
	firstStage := r.Stages[0]
	return timeutils.LocalDate(firstStage.StartDateTime.AsTime(), location)
}

func (r *Race) EndDateLocal(location *time.Location) time.Time {
	lastStage := r.Stages[len(r.Stages)-1]
	return timeutils.LocalDate(lastStage.StartDateTime.AsTime(), location)
}

func (r *Race) Result() []*ParticipantResultTime {
	var result []*ParticipantResultTime
	for _, stage := range r.Stages {
		generalTimeResults := stage.GeneralResults.GetTime()
		if len(generalTimeResults) > 0 {
			result = generalTimeResults
		}
	}
	return result
}

func (s *Stage) StartDateTimeLocal(location *time.Location) time.Time {
	return s.StartDateTime.AsTime().In(location)
}
