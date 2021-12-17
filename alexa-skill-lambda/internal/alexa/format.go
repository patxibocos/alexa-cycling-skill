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
		"uae-tour-2022":                         "El UAE Tour",
		"omloop-het-nieuwsblad-2022":            "La Omloop Het Nieuwsblad",
		"strade-bianche-2022":                   "La Strade Bianche",
		"paris-nice-2022":                       "La París Niza",
		"tirreno-adriatico-2022":                "El Tirreno Adriático",
		"milano-sanremo-2022":                   "La Milán San Remo",
		"volta-a-catalunya-2022":                "La Volta a Cataluña",
		"oxyclean-classic-brugge-de-panne-2022": "La Clásica Brujas La Panne",
		"e3-harelbeke-2022":                     "La Clásica E3 Saxo Bank",
		"gent-wevelgem-2022":                    "La Gante Wevelgem",
		"dwars-door-vlaanderen-2022":            "La Clásica A Través de Flandes",
		"ronde-van-vlaanderen-2022":             "El Tour de Flandes",
		"itzulia-basque-country-2022":           "La Vuelta al País Vasco",
		"amstel-gold-race-2022":                 "La Amstel Gold Race",
		"paris-roubaix-2022":                    "La París Roubaix",
		"la-fleche-wallone-2022":                "La Flecha Valona",
		"liege-bastogne-liege-2022":             "La Lieja Bastoña Lieja",
		"tour-de-romandie-2022":                 "El Tour de Romandía",
		"Eschborn-Frankfurt-2022":               "La Eschborn Frankfurt",
		"giro-d-italia-2022":                    "El Giro de Italia",
		"dauphine-2022":                         "El Critérium del Dauphiné",
		"tour-de-suisse-2022":                   "El Tour de Suiza",
		"tour-de-france-2022":                   "El Tour de Francia",
		"san-sebastian-2022":                    "La Clásica San Sebastián",
		"tour-de-pologne-2022":                  "El Tour de Polonia",
		"benelux-tour-2022":                     "El Tour del Benelux",
		"vuelta-a-espana-2022":                  "La Vuelta a España",
		"cyclassics-hamburg-2022":               "La Clásica Bemer",
		"bretagne-classic-2022":                 "La Clásica Bretaña",
		"gp-quebec-2022":                        "La Clásica Quebec",
		"gp-montreal-2022":                      "La Clásica Montreal",
		"il-lombardia-2022":                     "El Giro de Lombardía",
		"tour-of-guangxi-2022":                  "El Tour de Guangxi",
	}
)

func RiderFullName(rider *pcsscraper.Rider) string {
	return fmt.Sprintf("%s %s", rider.FirstName, rider.LastName)
}

func FormattedDate(time time.Time) string {
	return monthReplacer.Replace(time.Format("2 de January"))
}

func RaceName(raceID string) string {
	return raceIDToName[raceID]
}
