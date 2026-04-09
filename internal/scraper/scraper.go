package scraper

import (
	"fmt"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/models"
	"github.com/gocolly/colly"
)

var institutions []models.Institution

type InstitutionSource struct {
	URL  string
	Type constants.InstitutionType
}

var InstitutionRegistry = map[constants.InstitutionType]InstitutionSource{
	constants.FederalUniversity: {URL: constants.FederalUniversityURL, Type: constants.FederalUniversity},
	constants.StateUniversity:   {URL: constants.StateUniversityURL, Type: constants.StateUniversity},
	constants.PrivateUniversity: {URL: constants.PrivateUniversityURL, Type: constants.PrivateUniversity},

	// polytechnic
	constants.FederalPolytechnic: {URL: constants.FederalPolytechnicURL, Type: constants.FederalPolytechnic},
	constants.StatePolytechnic:   {URL: constants.StatePolytechnicURL, Type: constants.StatePolytechnic},
	constants.PrivatePolytechnic: {URL: constants.PrivatePolytechnicURL, Type: constants.PrivatePolytechnic},

	// college of education
	constants.FederalCollegeEduction: {URL: constants.FederalCollegeEductionURL, Type: constants.FederalCollegeEduction},
	constants.StateCollegeEduction:   {URL: constants.StateCollegeEductionURL, Type: constants.StateCollegeEduction},
	constants.PrivateCollegeEduction: {URL: constants.PrivateCollegeEductionURL, Type: constants.PrivateCollegeEduction},
}

type InstitutionScrapper struct {
	collector *colly.Collector
}

func NewInstitutionScrapper() *InstitutionScrapper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("www.nuc.edu.ng", "nuc.edu.ng", "education.gov.ng", "www.education.gov.ng"),
		colly.AllowURLRevisit(),
	)

	return &InstitutionScrapper{
		collector: c,
	}
}

func (s *InstitutionScrapper) ScrapeAllInstitution() {
	for _, institution := range InstitutionRegistry {
		_, _ = s.scrapeInstitution(institution.URL, string(institution.Type))
	}
}

func (s *InstitutionScrapper) scrapeInstitution(url, instituteType string) ([]model.Institution, error) {
	fmt.Println("Scraping", instituteType)

	s.collector.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		institutionName := e.ChildText(".column-2")
		institutionViceChancellor := e.ChildText(".column-3")
		institutionWebsite := e.ChildText(".column-4 a")
		institutionYearOfEstablishment := e.ChildText(".column-5")

		institution := models.Institution{
			Name:                institutionName,
			ViceChancellor:      institutionViceChancellor,
			YearOfEstablishment: institutionYearOfEstablishment,
			Url:                 institutionWebsite,
			Type:                instituteType,
		}

		institutions = append(institutions, institution)

	})

	fmt.Println("institutions", institutions)

	s.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	if err := s.collector.Visit(url); err != nil {
		return nil, fmt.Errorf("error visiting %s:%w", url, err)
	}
	return nil, nil
}
