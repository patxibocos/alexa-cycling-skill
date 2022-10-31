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
		"tour-down-under-2023":                  "El Santos Tour Down Under",
		"great-ocean-race-2023":                 "La Cadel Evans Great Ocean",
		"uae-tour-2023":                         "El U A E Tour",
		"omloop-het-nieuwsblad-2023":            "La Omloop Het Nieuwsblad",
		"strade-bianche-2023":                   "La Strade Bianche",
		"paris-nice-2023":                       "La París Niza",
		"tirreno-adriatico-2023":                "El Tirreno Adriático",
		"milano-sanremo-2023":                   "La Milán San Remo",
		"volta-a-catalunya-2023":                "La Volta a Cataluña",
		"oxyclean-classic-brugge-de-panne-2023": "La Clásica Brujas La Panne",
		"e3-harelbeke-2023":                     "La Clásica E3 Saxo Bank",
		"gent-wevelgem-2023":                    "La Gante Wevelgem",
		"dwars-door-vlaanderen-2023":            "La Clásica A Través de Flandes",
		"ronde-van-vlaanderen-2023":             "El Tour de Flandes",
		"itzulia-basque-country-2023":           "La Vuelta al País Vasco",
		"paris-roubaix-2023":                    "La París Rubé",
		"amstel-gold-race-2023":                 "La Amstel Gold Race",
		"la-fleche-wallone-2023":                "La Flecha Valona",
		"liege-bastogne-liege-2023":             "La Lieja Bastoña Lieja",
		"tour-de-romandie-2023":                 "El Tour de Romandía",
		"Eschborn-Frankfurt-2023":               "La Eschborn Frankfurt",
		"giro-d-italia-2023":                    "El Yiro de Italia",
		"dauphine-2023":                         "El Critérium del Dofiné",
		"tour-de-suisse-2023":                   "El Tour de Suiza",
		"tour-de-france-2023":                   "El Tour de Francia",
		"san-sebastian-2023":                    "La Clásica San Sebastián",
		"tour-de-pologne-2023":                  "El Tour de Polonia",
		"cyclassics-hamburg-2023":               "La Clásica Bemer",
		"benelux-tour-2023":                     "El Tour del Benelux",
		"vuelta-a-espana-2023":                  "La Vuelta a España",
		"bretagne-classic-2023":                 "La Clásica Bretaña",
		"gp-quebec-2023":                        "La Clásica Quebec",
		"gp-montreal-2023":                      "La Clásica Montreal",
		"il-lombardia-2023":                     "El Giro de Lombardía",
		"tour-of-guangxi-2023":                  "El Tour de Guangxi",
	}
)

func riderFullName(rider *pcsscraper.Rider) string {
	return fmt.Sprintf("%s %s", rider.FirstName, rider.LastName)
}

func formattedDate(time time.Time) string {
	return monthReplacer.Replace(time.Format("2 de January"))
}

func raceName(raceID string) string {
	return raceIDToName[raceID]
}
