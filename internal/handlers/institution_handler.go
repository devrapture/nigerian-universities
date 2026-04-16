package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
)

type Institution struct {
	Name                string `json:"name"`
	ViceChancellor      string `json:"vice_chancellor"`
	YearOfEstablishment string `json:"year_of_establishment"`
	Type                string `json:"type"`
	Url                 string `json:"url"`
}

type InstitutionHandler struct {
	institutionService service.InstitutionService
}

func NewInstitutionHandler(institutionService service.InstitutionService) *InstitutionHandler {
	return &InstitutionHandler{
		institutionService: institutionService,
	}
}

func (h *InstitutionHandler) GetAllInstitutions(c *gin.Context) {
	queryDTO, err := parseListQuery(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	allInstitution, total, err := h.institutionService.GetAllInstitutions(c.Request.Context(), queryDTO)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
		return
	}

	meta := &utils.PaginationMeta{
		Page:    queryDTO.Page,
		PerPage: queryDTO.Limit,
		Total:   total,
		Pages:   int64(math.Ceil(float64(total) / float64(queryDTO.Limit))),
	}
	utils.SuccessResponse(c, http.StatusOK, "fetched all institutions", allInstitution, meta)

}

// parseListQuery manually parses query params to give clearer error messages than the default binder.
func parseListQuery(c *gin.Context) (dto.ListInstitutionQuery, error) {
	const maxLimit = 100
	q := dto.ListInstitutionQuery{
		Search: strings.TrimSpace(c.Query("search")),
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return q, friendlyNumErr("page", pageStr)
	}
	q.Page = page
	if page < 1 {
		return q, fmt.Errorf("query parameter 'page' must be at least 1, got '%d'", page)
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return q, friendlyNumErr("limit", limitStr)
	}
	if limit < 1 {
		return q, fmt.Errorf("query parameter 'limit' must be at least 1, got '%d'", limit)
	}
	if limit > maxLimit {
		return q, fmt.Errorf("query parameter 'limit' must be <= %d", maxLimit)
	}

	q.Limit = limit

	// optional type; leave empty if not provided
	if t := strings.TrimSpace(c.Query("type")); t != "" {
		q.Type = constants.InstitutionType(t)
		if !isValidInstitutionType(q.Type) {
			return q, fmt.Errorf("query parameter 'type' must be one of [%s]", strings.Join(validInstitutionTypes(), ", "))
		}
	}

	return q, nil
}

func friendlyNumErr(field, value string) error {
	return fmt.Errorf("query parameter '%s' must be a number, got '%s'", field, value)
}

func isValidInstitutionType(t constants.InstitutionType) bool {
	for _, v := range validInstitutionTypes() {
		if string(t) == v {
			return true
		}
	}
	return false
}

func validInstitutionTypes() []string {
	return []string{
		string(constants.FederalUniversity),
		string(constants.StateUniversity),
		string(constants.PrivateUniversity),
		string(constants.FederalPolytechnic),
		string(constants.StatePolytechnic),
		string(constants.PrivatePolytechnic),
		string(constants.FederalCollegeEducation),
		string(constants.StateCollegeEducation),
		string(constants.PrivateCollegeEducation),
	}
}
