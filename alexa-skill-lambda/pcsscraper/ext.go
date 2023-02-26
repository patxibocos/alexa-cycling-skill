package pcsscraper

import (
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/timeutils"
	"time"
)

func (r *Race) StartDateLocal(location *time.Location) time.Time {
	return timeutils.LocalDate(r.StartDate.AsTime(), location)
}

func (r *Race) EndDateLocal(location *time.Location) time.Time {
	return timeutils.LocalDate(r.EndDate.AsTime(), location)
}

func (s *Stage) StartDateTimeLocal(location *time.Location) time.Time {
	return s.StartDateTime.AsTime().In(location)
}
