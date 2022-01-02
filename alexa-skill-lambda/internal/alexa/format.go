package alexa

import (
	"fmt"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/internal/cycling"
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
		"uae-tour-2022":                         "El U A E Tour",
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
		"paris-roubaix-2022":                    "La París Rubé",
		"la-fleche-wallone-2022":                "La Flecha Valona",
		"liege-bastogne-liege-2022":             "La Lieja Bastoña Lieja",
		"tour-de-romandie-2022":                 "El Tour de Romandía",
		"Eschborn-Frankfurt-2022":               "La Eschborn Frankfurt",
		"giro-d-italia-2022":                    "El Yiro de Italia",
		"dauphine-2022":                         "El Critérium del Dofiné",
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

	stageTypeToText = map[pcsscraper.Stage_Type]string{
		pcsscraper.Stage_TYPE_FLAT:                    "llano",
		pcsscraper.Stage_TYPE_HILLS_FLAT_FINISH:       "de media montaña con final llano",
		pcsscraper.Stage_TYPE_HILLS_UPHILL_FINISH:     "de media montaña con final en alto",
		pcsscraper.Stage_TYPE_MOUNTAINS_FLAT_FINISH:   "de montaña con final llano",
		pcsscraper.Stage_TYPE_MOUNTAINS_UPHILL_FINISH: "de montaña con final en alto",
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

func stageType(stageType pcsscraper.Stage_Type) string {
	return stageTypeToText[stageType]
}

func phraseWithTop3(phrase string, top3 *cycling.Top3) string {
	return fmt.Sprintf(
		phrase,
		riderFullName(top3.First.Rider),
		riderFullName(top3.Second.Rider),
		riderFullName(top3.Third.Rider),
	)
}

func singularOrPlural(word string, amount int64) string {
	if amount == 1 {
		return word
	}
	return word + "s"
}

func getGapMessage(gap int64) string {
	if gap == 0 {
		return "con el mismo tiempo"
	}
	if gap < 60 {
		return fmt.Sprintf("a %d %s", gap, singularOrPlural("segundo", gap))
	}
	minutes := gap / 60
	seconds := gap % 60
	message := fmt.Sprintf("a %d %s", minutes, singularOrPlural("minuto", minutes))
	if seconds > 0 {
		message += fmt.Sprintf(" y %d %s", seconds, singularOrPlural("segundo", seconds))
	}
	return message
}

func phraseWithTop3AndGaps(phrase string, top3 *cycling.Top3) string {
	firstToSecondGap := top3.Second.Time - top3.First.Time
	secondToThirdGap := top3.Third.Time - top3.Second.Time
	firstToSecondMessage := getGapMessage(firstToSecondGap)
	secondToThirdMessage := getGapMessage(secondToThirdGap)
	return fmt.Sprintf(
		phrase,
		riderFullName(top3.First.Rider),
		firstToSecondMessage,
		riderFullName(top3.Second.Rider),
		secondToThirdMessage,
		riderFullName(top3.Third.Rider),
	)
}
