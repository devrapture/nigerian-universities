package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"gorm.io/gorm"
)

type InstitutionRepository interface {
	UpsertMany(ctx context.Context, institutions []model.Institution) error
	FindAll(ctx context.Context) ([]model.Institution, error)
}

type institutionRepository struct {
	db *gorm.DB
}

func NewInstitutionRepository(db *gorm.DB) InstitutionRepository {
	return &institutionRepository{
		db: db,
	}
}

func (r *institutionRepository) UpsertMany(ctx context.Context, institutions []model.Institution) error {
	for _, institution := range institutions {
		fmt.Println("upserting", institution.Name)
		var existing model.Institution
		err := r.db.WithContext(ctx).Where("name = ? AND type = ?", institution.Name, institution.Type).First(&existing).Error
		if err == nil {
			now := time.Now()
			existing.ViceChancellor = institution.ViceChancellor
			existing.Website = institution.Website
			existing.YearOfEstablishment = institution.YearOfEstablishment
			existing.LastScrapedAt = &now

			if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
				return err
			}

			continue
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			institution.LastScrapedAt = &now
			if err := r.db.WithContext(ctx).Create(&institution).Error; err != nil {
				return err

			}
			continue
		}

		return err
	}

	return nil
}

func (r *institutionRepository) FindAll(ctx context.Context) ([]model.Institution, error) {
	var institutions []model.Institution
	if err := r.db.WithContext(ctx).Find(&institutions).Error; err != nil {
		return nil, err
	}
	return institutions, nil
}
