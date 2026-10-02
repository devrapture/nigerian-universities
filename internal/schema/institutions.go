package schema

import "time"

type InstitutionItem struct {
	ID                  string     `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                string     `json:"name" example:"University of Lagos"`
	Type                string     `json:"type" example:"federal-university"`
	ViceChancellor      string     `json:"vice_chancellor" example:"Prof. Olatunji Afolabi Oyelana"`
	Website             string     `json:"website" example:"https://www.unilag.edu.ng"`
	YearOfEstablishment string     `json:"year_of_establishment" example:"1962"`
	LastScrapedAt       *time.Time `json:"last_scraped_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type InstitutionListResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    []InstitutionItem `json:"data"`
	Meta    PaginationMeta    `json:"meta,omitempty"`
}

type InstitutionUnauthorizedError struct {
	Code    string `json:"code" example:"INVALID_API_KEY"`
	Message string `json:"message" example:"Invalid or missing X-API-Key header"`
}

type InstitutionUnauthorizedResponse struct {
	Success bool                          `json:"success" example:"false"`
	Error   *InstitutionUnauthorizedError `json:"error,omitempty"`
}

type InstitutionBadRequestError struct {
	Code    string `json:"code" example:"BAD_REQUEST"`
	Message string `json:"message" example:"query parameter 'page' must be a number, got 'abc'"`
}

type InstitutionBadRequestResponse struct {
	Success bool                        `json:"success" example:"false"`
	Error   *InstitutionBadRequestError `json:"error,omitempty"`
}

type InstitutionInternalServerError struct {
	Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"An unexpected error occurred"`
}

type InstitutionInternalServerErrorResponse struct {
	Success bool                            `json:"success" example:"false"`
	Error   *InstitutionInternalServerError `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
	Pages   int64 `json:"pages"`
}

type InstitutionFieldChange struct {
	Field    string `json:"field" example:"website"`
	OldValue string `json:"old_value" example:"https://old.example.edu.ng"`
	NewValue string `json:"new_value" example:"https://new.example.edu.ng"`
}

type ChangeLogItem struct {
	ID              string                   `json:"id"`
	InstitutionID   string                   `json:"institution_id"`
	InstitutionName string                   `json:"institution_name" example:"University of Lagos"`
	InstitutionType string                   `json:"institution_type" example:"federal-university"`
	ChangeType      string                   `json:"change_type" example:"updated"`
	Changes         []InstitutionFieldChange `json:"changes"`
	ChangedAt       time.Time                `json:"changed_at"`
}

type ChangeLogListResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    []ChangeLogItem `json:"data"`
	Meta    PaginationMeta  `json:"meta"`
}

type ScrapeTriggerData struct {
	Status string `json:"status" example:"completed"`
	Count  int    `json:"count,omitempty" example:"120"`
}

type ScrapeTriggerResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"scrape completed"`
	Data    ScrapeTriggerData `json:"data"`
}

type ScrapeTimeoutError struct {
	Code    string `json:"code" example:"TIMEOUT"`
	Message string `json:"message" example:"scrape timed out"`
}

type ScrapeTimeoutResponse struct {
	Success bool                `json:"success" example:"false"`
	Error   *ScrapeTimeoutError `json:"error,omitempty"`
}
