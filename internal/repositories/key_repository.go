package repositories

import (
	"context"
	"time"

	apperrors "github.com/coolpythoncodes/nigerian-universities/internal/errors"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeyRepository interface {
	CreateKey(ctx context.Context, userID uuid.UUID, rawKey string) (*model.ProductKey, error)
	GetAllKeys(ctx context.Context, userID uuid.UUID) ([]model.ProductKey, error)
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
	key := &model.ProductKey{
		UserID:    userID,
		KeyHash:   utils.HashKey(rawKey),
		KeyPrefix: rawKey[:8],
	}

	if err := r.db.WithContext(ctx).Create(key).Error; err != nil {
		return nil, err
	}
	return key, nil
}

func (r *keyRepository) GetAllKeys(ctx context.Context, userID uuid.UUID) ([]model.ProductKey, error) {
	var keys []model.ProductKey

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		return nil, err
	}

	return keys, nil
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
