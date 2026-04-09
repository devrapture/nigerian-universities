package main

import "github.com/coolpythoncodes/nigerian-universities/internal/scraper"

func main() {
	s:= scraper.NewInstitutionScrapper()

	s.ScrapeAllInstitution()
}
