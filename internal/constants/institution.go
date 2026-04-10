package constants

// InstitutionType represents the category of an institution.
type InstitutionType string

const (
	// Universities
	FederalUniversity InstitutionType = "federal-university"
	StateUniversity   InstitutionType = "state-university"
	PrivateUniversity InstitutionType = "private-university"

	// Polytechnics
	FederalPolytechnic InstitutionType = "federal-polytechnic"
	StatePolytechnic   InstitutionType = "state-polytechnic"
	PrivatePolytechnic InstitutionType = "private-polytechnic"

	// Colleges of education
	FederalCollegeEducation InstitutionType = "federal-college-education"
	StateCollegeEducation   InstitutionType = "state-college-education"
	PrivateCollegeEducation InstitutionType = "private-college-education"
)

// URLs for scraping each institution type.
const (
	FederalUniversityURL = "https://www.nuc.edu.ng/nigerian-univerisities/federal-univeristies/"
	StateUniversityURL   = "https://www.nuc.edu.ng/nigerian-univerisities/state-univerisity/"
	PrivateUniversityURL = "https://www.nuc.edu.ng/nigerian-univerisities/private-univeristies/"

	// polytechnic
	FederalPolytechnicURL = "https://education.gov.ng/government-polytechnics/"
	StatePolytechnicURL   = "https://education.gov.ng/state-polytechnics/"
	PrivatePolytechnicURL = "https://education.gov.ng/private-polytechnics/"

	// college of education
	FederalCollegeEducationURL = "https://education.gov.ng/federal-college-of-education/"
	StateCollegeEducationURL   = "https://education.gov.ng/state-college-of-education/"
	PrivateCollegeEducationURL = "https://education.gov.ng/private-college-of-education/"
)
