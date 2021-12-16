package cycling

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
)

func RiderFullName(rider *pcsscraper.Rider) string {
	return fmt.Sprintf("%s %s", rider.FirstName, rider.LastName)
}
