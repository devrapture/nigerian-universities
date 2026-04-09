package scraper

import (
	"fmt"

	"github.com/coolpythoncodes/nigerian-universities/models"
	"github.com/coolpythoncodes/nigerian-universities/utils"
	"github.com/gocolly/colly"
)

var universities []models.Institution

func ScrapeUniversities() {
	pagesToScrape := map[string]string{
		"https://www.nuc.edu.ng/nigerian-univerisities/federal-univeristies/": "Federal",
		"https://www.nuc.edu.ng/nigerian-univerisities/state-univerisity/":    "State",
		"https://www.nuc.edu.ng/nigerian-univerisities/private-univeristies/": "Private",
	}
	c := colly.NewCollector(
		colly.AllowedDomains("www.nuc.edu.ng", "nuc.edu.ng"),
	)

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		universityType := pagesToScrape[e.Request.URL.String()]

		nameOfUniversity := e.ChildText(".column-2")
		nameOfChancellor := e.ChildText(".column-3")
		universityWebsite := e.ChildText(".column-4 a")
		yearOfEstablishment := e.ChildText(".column-5")

		university := models.Institution{
			Name:                nameOfUniversity,
			ViceChancellor:      nameOfChancellor,
			YearOfEstablishment: yearOfEstablishment,
			Url:                 universityWebsite,
			Type:                universityType,
		}

		universities = append(universities, university)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for url := range pagesToScrape {
		c.Visit(url)
	}
	utils.WriteJSON("nigerian-university.json", universities)
}
