package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	apperrors "github.com/coolpythoncodes/nigerian-universities/internal/errors"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeyRepository interface {
	CreateKey(ctx context.Context, userID uuid.UUID, rawKey string) (*model.ProductKey, error)
	GetAllKeys(ctx context.Context, userID uuid.UUID, queryDTO dto.ListInstitutionQuery) ([]model.ProductKey, int64, error)
	RevokeKey(ctx context.Context, userID, keyID uuid.UUID) error
}

type keyRepository struct {
	db *gorm.DB
}

func NewKeyRepository(DB *gorm.DB) KeyRepository {
	return &keyRepository{
		db: DB,
	}
}

func (r *keyRepository) CreateKey(ctx context.Context, userID uuid.UUID, rawKey string) (*model.ProductKey, error) {
	now := time.Now().UTC()
	key := &model.ProductKey{
		UserID:    userID,
		KeyHash:   utils.HashKey(rawKey),
		KeyPrefix: rawKey[:8],
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(&model.ProductKey{}).Where("user_id = ? AND is_active = true", userID).Updates(map[string]interface{}{
			"is_active":  false,
			"revoked_at": now,
		}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Create(key).Error; err != nil {
			return fmt.Errorf("failed to create key")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (r *keyRepository) GetAllKeys(ctx context.Context, userID uuid.UUID, queryDTO dto.ListInstitutionQuery) ([]model.ProductKey, int64, error) {
	var total int64
	var keys []model.ProductKey

	query := r.db.WithContext(ctx).Model(&model.ProductKey{}).Where("user_id = ?", userID)

	query.Count(&total)
	query = query.Order("created_at DESC").Offset((queryDTO.Page - 1) * queryDTO.Limit).Limit(queryDTO.Limit)
	if err := query.Find(&keys).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get keys")
	}
	return keys, total, nil
}

func (r *keyRepository) RevokeKey(ctx context.Context, userID, keyID uuid.UUID) error {
	var keyToRevoke model.ProductKey

	if err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, keyID).First(&keyToRevoke).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrKeyNotFound
		}
		return err
	}

	if userID != keyToRevoke.UserID {
		return apperrors.ErrUnauthorized
	}

	now := time.Now().UTC()

	keyToRevoke.IsActive = false
	keyToRevoke.RevokedAt = &now

	if err := r.db.WithContext(ctx).Save(&keyToRevoke).Error; err != nil {
		return err
	}

	return nil
}
