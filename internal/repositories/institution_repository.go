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
	FindChangeLogs(ctx context.Context, queryDto dto.ListChangeLogQuery) ([]model.ChangeLog, int64, error)
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, institution := range institutions {
			fmt.Println("upserting", institution.Name)
			var existing model.Institution
			err := tx.Where("name = ? AND type = ?", institution.Name, institution.Type).First(&existing).Error
			now := time.Now()

			if err == nil {
				changes := institutionChanges(existing, institution)
				existing.ViceChancellor = institution.ViceChancellor
				existing.Website = institution.Website
				existing.YearOfEstablishment = institution.YearOfEstablishment
				existing.LastScrapedAt = &now

				if err := tx.Save(&existing).Error; err != nil {
					return err
				}

				if len(changes) > 0 {
					if err := tx.Create(newChangeLog(existing, "updated", changes, now)).Error; err != nil {
						return err
					}
				}
				continue
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				institution.LastScrapedAt = &now
				if err := tx.Create(&institution).Error; err != nil {
					return err
				}
				if err := tx.Create(newChangeLog(institution, "created", model.InstitutionFieldChanges{}, now)).Error; err != nil {
					return err
				}
				continue
			}

			return err
		}

		return nil
	})
}

func institutionChanges(existing, incoming model.Institution) model.InstitutionFieldChanges {
	changes := make(model.InstitutionFieldChanges, 0, 3)
	if existing.ViceChancellor != incoming.ViceChancellor {
		changes = append(changes, model.InstitutionFieldChange{Field: "vice_chancellor", OldValue: existing.ViceChancellor, NewValue: incoming.ViceChancellor})
	}
	if existing.Website != incoming.Website {
		changes = append(changes, model.InstitutionFieldChange{Field: "website", OldValue: existing.Website, NewValue: incoming.Website})
	}
	if existing.YearOfEstablishment != incoming.YearOfEstablishment {
		changes = append(changes, model.InstitutionFieldChange{Field: "year_of_establishment", OldValue: existing.YearOfEstablishment, NewValue: incoming.YearOfEstablishment})
	}
	return changes
}

func newChangeLog(institution model.Institution, changeType string, changes model.InstitutionFieldChanges, changedAt time.Time) *model.ChangeLog {
	return &model.ChangeLog{
		InstitutionID:   institution.ID,
		InstitutionName: institution.Name,
		InstitutionType: institution.Type,
		ChangeType:      changeType,
		Changes:         changes,
		ChangedAt:       changedAt,
	}
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

func (r *institutionRepository) FindChangeLogs(ctx context.Context, queryDto dto.ListChangeLogQuery) ([]model.ChangeLog, int64, error) {
	var changeLogs []model.ChangeLog
	var total int64
	query := r.db.WithContext(ctx).Model(&model.ChangeLog{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("changed_at DESC").
		Order("id DESC").
		Offset((queryDto.Page - 1) * queryDto.Limit).
		Limit(queryDto.Limit).
		Find(&changeLogs).Error; err != nil {
		return nil, 0, err
	}

	return changeLogs, total, nil
}
