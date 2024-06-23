package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"strings"
	"time"
)

var (
	monthReplacer = strings.NewReplacer(
		"January", "Enero",
		"February", "Febrero",
		"March", "Marzo",
		"April", "Abril",
		"May", "Mayo",
		"June", "Junio",
		"July", "Julio",
		"August", "Agosto",
		"September", "Septiembre",
		"October", "Octubre",
		"November", "Noviembre",
		"December", "Diciembre")

	raceIDToName = map[string]string{
		"tour-down-under":                  "El Santos Tour Down Under",
		"great-ocean-race":                 "La Cadel Evans Great Ocean",
		"uae-tour":                         "El U A E Tour",
		"omloop-het-nieuwsblad":            "La Omloop Het Nieuwsblad",
		"strade-bianche":                   "La Strade Bianche",
		"paris-nice":                       "La París Niza",
		"tirreno-adriatico":                "El Tirreno Adriático",
		"milano-sanremo":                   "La Milán San Remo",
		"volta-a-catalunya":                "La Volta a Cataluña",
		"oxyclean-classic-brugge-de-panne": "La Clásica Brujas La Panne",
		"e3-harelbeke":                     "La Clásica E3 Saxo Bank",
		"gent-wevelgem":                    "La Gante Wevelgem",
		"dwars-door-vlaanderen":            "La Clásica A Través de Flandes",
		"ronde-van-vlaanderen":             "El Tour de Flandes",
		"itzulia-basque-country":           "La Vuelta al País Vasco",
		"paris-roubaix":                    "La París Rubé",
		"amstel-gold-race":                 "La Amstel Gold Race",
		"la-fleche-wallone":                "La Flecha Valona",
		"liege-bastogne-liege":             "La Lieja Bastoña Lieja",
		"tour-de-romandie":                 "El Tour de Romandía",
		"eschborn-frankfurt":               "La Eschborn Frankfurt",
		"giro-d-italia":                    "El Yiro de Italia",
		"dauphine":                         "El Critérium del Dofiné",
		"tour-de-suisse":                   "El Tour de Suiza",
		"tour-de-france":                   "El Tour de Francia",
		"san-sebastian":                    "La Clásica San Sebastián",
		"tour-de-pologne":                  "El Tour de Polonia",
		"cyclassics-hamburg":               "La Clásica Bemer",
		"renewi-tour":                      "El Tour de Renewi",
		"benelux-tour":                     "El Tour del Benelux",
		"vuelta-a-espana":                  "La Vuelta a España",
		"bretagne-classic":                 "La Clásica Bretaña",
		"gp-quebec":                        "La Clásica Quebec",
		"gp-montreal":                      "La Clásica Montreal",
		"il-lombardia":                     "El Giro de Lombardía",
		"tour-of-guangxi":                  "El Tour de Guangxi",
	}
)

func riderFullName(rider *pcsscraper.Rider) string {
	return fmt.Sprintf("%s %s", rider.FirstName, rider.LastName)
}

func formattedDate(time time.Time) string {
	return monthReplacer.Replace(time.Format("2 de January"))
}

func nameForRace(race *pcsscraper.Race) string {
	// raceID contains -YYYY, so we remove it to calculate map's key
	lastDashIndex := strings.LastIndex(race.Id, "-")
	mapKey := strings.ToLower(race.Id[:lastDashIndex])
	return raceIDToName[mapKey]
}
