package model

import (
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Institution struct {
	ID            uuid.UUID                 `json:"id" gorm:"type:uuid;primary_key" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name          string                    `json:"name" gorm:"type:text;not null;index:idx_universities_name" example:"University of Lagos"`
	Type          constants.InstitutionType `json:"type" gorm:"type:text;not null;index:idx_universities_type" example:"Federal"`
	LastScrapedAt *time.Time                `json:"last_scraped_at,omitempty" example:"2021-01-01T00:00:00Z"`
	CreatedAt     time.Time                 `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time                 `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt     gorm.DeletedAt            `json:"-" gorm:"index" example:"2021-01-01T00:00:00Z"`
}

// BeforeCreate generates UUID before inserting
func (u *Institution) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
