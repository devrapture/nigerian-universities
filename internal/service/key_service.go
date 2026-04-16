package service

import (
	"context"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/google/uuid"
)

type KeyService interface {
	HandleCreateKey(ctx context.Context, userID uuid.UUID) (*model.ProductKey, string, error)
	HandleGetAllKeys(ctx context.Context, userID uuid.UUID,page, perPage int) ([]model.ProductKey, int64, error)
	HandleRevokeKey(ctx context.Context, userID, keyID uuid.UUID) error
}

type keyService struct {
	keyRepo repositories.KeyRepository
}

func NewKeyService(keyRepo repositories.KeyRepository) KeyService {
	return &keyService{
		keyRepo: keyRepo,
	}
}

func (s *keyService) HandleCreateKey(ctx context.Context, userID uuid.UUID) (*model.ProductKey, string, error) {
	rawKey, err := utils.GenerateRawKey()
	if err != nil {
		return nil, "", err
	}
	productKey, err := s.keyRepo.CreateKey(ctx, userID, rawKey)

	if err != nil {
		return nil, "", err
	}
	return productKey, rawKey, nil
}

func (s *keyService) HandleGetAllKeys(ctx context.Context, userID uuid.UUID,page, perPage int) ([]model.ProductKey, int64, error) {
	return s.keyRepo.GetAllKeys(ctx, userID, page, perPage)
}

func (s *keyService) HandleRevokeKey(ctx context.Context, userID, keyID uuid.UUID) error {
	return s.keyRepo.RevokeKey(ctx, userID, keyID)
}
