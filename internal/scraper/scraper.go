package scraper

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

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
	constants.FederalCollegeEducation: {URL: constants.FederalCollegeEducationURL, Type: constants.FederalCollegeEducation},
	constants.StateCollegeEducation:   {URL: constants.StateCollegeEducationURL, Type: constants.StateCollegeEducation},
	constants.PrivateCollegeEducation: {URL: constants.PrivateCollegeEducationURL, Type: constants.PrivateCollegeEducation},
}

type InstitutionScrapper struct{}

func NewInstitutionScrapper() *InstitutionScrapper {
	return &InstitutionScrapper{}
}

func (s *InstitutionScrapper) ScrapeAllInstitution(ctx context.Context) ([]model.Institution, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	sources := make([]InstitutionSource, 0, len(InstitutionRegistry))
	for _, source := range InstitutionRegistry {
		sources = append(sources, source)
	}

	out := make([][]model.Institution, len(sources))
	g, ctx := errgroup.WithContext(ctx)
	for i, source := range sources {
		i, source := i, source
		g.Go(func() error {
			if err := ctx.Err(); err != nil {
				return err
			}
			institutions, err := s.scrapeInstitution(ctx, source.URL, source.Type)
			if err != nil {
				return err
			}
			fmt.Println("scraped institutions", institutions)
			out[i] = institutions
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	var allInstitutions []model.Institution
	for _, institutions := range out {
		allInstitutions = append(allInstitutions, institutions...)
	}
	fmt.Println("all institutions", allInstitutions)
	return allInstitutions, nil
}

func remainingTimeout(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 60 * time.Second
	}
	remaining := time.Until(deadline)
	if remaining < time.Millisecond {
		return time.Millisecond
	}
	return remaining
}

func (s *InstitutionScrapper) scrapeInstitution(ctx context.Context, url string, instituteType constants.InstitutionType) ([]model.Institution, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// use a fresh slice per scrape to avoid cross-source accumulation
	institutions := make([]model.Institution, 0, 64)

	timeout := remainingTimeout(ctx)
	headerTimeout := timeout
	if headerTimeout > 20*time.Second {
		headerTimeout = 20 * time.Second
	}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("www.nuc.edu.ng", "nuc.edu.ng", "education.gov.ng", "www.education.gov.ng"),
		colly.AllowURLRevisit(),
	)

	// education.gov.ng pages are slow; cap timeouts to the remaining request deadline
	collector.WithTransport(&http.Transport{ResponseHeaderTimeout: headerTimeout})
	collector.SetRequestTimeout(timeout)

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
		if strings.Contains(string(instituteType), "college-education") {
			institutionViceChancellor = ""
		}

		institution := model.Institution{
			Name:                institutionName,
			ViceChancellor:      institutionViceChancellor,
			YearOfEstablishment: institutionYearOfEstablishment,
			Website:             institutionWebsite,
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

		institution := model.Institution{
			Name:                institutionName,
			ViceChancellor:      "",
			YearOfEstablishment: institutionYearOfEstablishment,
			Website:             institutionWebsite,
			Type:                instituteType,
		}

		institutions = append(institutions, institution)
	})

	// fmt.Println("institutions", institutions)

	collector.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			return
		}
		fmt.Println("Visiting", r.URL.String())
	})
	if err := collector.Visit(url); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, fmt.Errorf("error visiting %s: %w", url, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return institutions, nil
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
