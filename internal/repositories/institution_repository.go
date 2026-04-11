package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"gorm.io/gorm"
)

type InstitutionRepository interface {
	UpsertMany(ctx context.Context, institutions []model.Institution) error
	FindAll(ctx context.Context, queryDto dto.ListInstitutionQuery) ([]model.Institution, int64, error)
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

func (r *institutionRepository) FindAll(ctx context.Context, queryDto dto.ListInstitutionQuery) ([]model.Institution, int64, error) {
	var institutions []model.Institution
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Institution{})

	if queryDto.Type != "" {
		query = query.Where("type=?", queryDto.Type)
	}

	if queryDto.Search != "" {
		search := strings.ToLower(queryDto.Search)
		query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	query = query.Order("name ASC").
		Offset((queryDto.Page - 1) * queryDto.Limit).
		Limit(queryDto.Limit)

	if err := query.Find(&institutions).Error; err != nil {
		return nil, 0, err
	}

	return institutions, total, nil
}
