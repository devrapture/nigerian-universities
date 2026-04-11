package dto

import "github.com/coolpythoncodes/nigerian-universities/internal/constants"

type ListInstitutionQuery struct {
	Page   int                       `form:"page" default:"1"`
	Limit  int                       `form:"limit" default:"10"`
	Type   constants.InstitutionType `form:"type" example:"federal-university"`
	Search string                    `form:"search" example:"university of lagos"`
}
