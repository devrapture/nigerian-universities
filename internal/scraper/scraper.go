package scraper

import (
	"fmt"
	"strings"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/models"
	"github.com/coolpythoncodes/nigerian-universities/utils"
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

type InstitutionScrapper struct{}

func NewInstitutionScrapper() *InstitutionScrapper {
	return &InstitutionScrapper{}
}

func (s *InstitutionScrapper) ScrapeAllInstitution() {
	for _, institution := range InstitutionRegistry {
		_, _ = s.scrapeInstitution(institution.URL, string(institution.Type))
	}
}

func (s *InstitutionScrapper) scrapeInstitution(url, instituteType string) ([]model.Institution, error) {
	fmt.Println("Scraping", instituteType)

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("www.nuc.edu.ng", "nuc.edu.ng", "education.gov.ng", "www.education.gov.ng"),
		colly.AllowURLRevisit(),
	)

	// for nuc
	collector.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		// ignore rows that belong to the education.gov.ng table (handled below)
		if e.DOM.ParentsFiltered("table#tablepress-19").Length() > 0 {
			return
		}

		institutionName := strings.TrimSpace(e.ChildText(".column-2"))
		institutionViceChancellor := strings.TrimSpace(e.ChildText(".column-3"))
		institutionWebsite := strings.TrimSpace(firstNonEmpty(
			e.ChildAttr(".column-4 a", "href"),
			e.ChildText(".column-4"),
		))
		institutionYearOfEstablishment := strings.TrimSpace(e.ChildText(".column-5"))

		if institutionName == "" {
			return
		}

		// colleges of education don't have vice chancellors
		if strings.Contains(instituteType, "college-education") {
			institutionViceChancellor = ""
		}

		institution := models.Institution{
			Name:                institutionName,
			ViceChancellor:      institutionViceChancellor,
			YearOfEstablishment: institutionYearOfEstablishment,
			Url:                 institutionWebsite,
			Type:                instituteType,
		}

		institutions = append(institutions, institution)

	})

	// for education.gov.ng (tablepress layout)
	collector.OnHTML("table#tablepress-19 tbody tr", func(e *colly.HTMLElement) {
		institutionName := strings.TrimSpace(e.ChildText(".column-2"))
		institutionWebsite := strings.TrimSpace(firstNonEmpty(
			e.ChildAttr(".column-3 a", "href"),
			e.ChildText(".column-3"),
		))
		institutionYearOfEstablishment := strings.TrimSpace(e.ChildText(".column-4"))

		// skip rows without a name (they are usually secondary link rows)
		if institutionName == "" {
			return
		}

		institution := models.Institution{
			Name:                institutionName,
			ViceChancellor:      "",
			YearOfEstablishment: institutionYearOfEstablishment,
			Url:                 institutionWebsite,
			Type:                instituteType,
		}

		institutions = append(institutions, institution)

	})

	fmt.Println("institutions", institutions)
	utils.WriteJSON("nigerian-institutions.json", institutions)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	if err := collector.Visit(url); err != nil {
		return nil, fmt.Errorf("error visiting %s:%w", url, err)
	}
	return nil, nil
}

// firstNonEmpty returns the first string that is not empty after trimming.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
