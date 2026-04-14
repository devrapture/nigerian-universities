package model

import (
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email       string         `json:"email" gorm:"uniqueIndex:idx_users_email;not null" example:"john.doe@example.com" description:"User's email"`
	Name        string         `json:"name" gorm:"type:text;not null" example:"John Doe" description:"User's name"`
	AvatarURL   string         `json:"avatar_url" gorm:"type:text;not null" example:"https://example.com/avatar.png" description:"User's avatar URL"`
	Provider    string         `json:"provider" gorm:"type:text;not null" example:"google" description:"User's provider"`
	ProviderID  string         `json:"provider_id" gorm:"type:text;not null" example:"1234567890" description:"User's provider ID"`
	ProductKeys []ProductKey   `json:"product_keys,omitempty" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time      `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index" example:"2021-01-01T00:00:00Z"`
}

type Institution struct {
	ID                  uuid.UUID                 `json:"id" gorm:"type:uuid;primary_key" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                string                    `json:"name" gorm:"type:text;not null;index:idx_universities_name" example:"University of Lagos"`
	Type                constants.InstitutionType `json:"type" gorm:"type:text;not null;index:idx_universities_type" example:"Federal"`
	ViceChancellor      string                    `json:"vice_chancellor" gorm:"type:text;not null;index:idx_universities_vice_chancellor" example:"Prof. Olatunji Afolabi Oyelana"`
	Website             string                    `json:"website" gorm:"type:text;not null;index:idx_universities_website" example:"https://www.aun.edu.ng"`
	YearOfEstablishment string                    `json:"year_of_establishment" gorm:"type:text;not null;index:idx_universities_year_of_establishment" example:"1999"`
	LastScrapedAt       *time.Time                `json:"last_scraped_at,omitempty" example:"2021-01-01T00:00:00Z"`
	CreatedAt           time.Time                 `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt           time.Time                 `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt           gorm.DeletedAt            `json:"-" gorm:"index" example:"2021-01-01T00:00:00Z"`
}

type ProductKey struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index:idx_product_keys_user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	KeyHash   string         `json:"-" gorm:"not null;uniqueIndex" example:"123e4567-e89b-12d3-a456-426614174000"`
	KeyPrefix string         `json:"key_prefix" gorm:"not null"`
	IsActive  bool           `json:"is_active" gorm:"not null;default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" example:"2021-01-01T00:00:00Z"`
}

// BeforeCreate generates UUID before inserting
func (u *Institution) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate generates UUID before inserting
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate generates UUID before inserting
func (p *ProductKey) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
