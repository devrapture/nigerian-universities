package service

import (
	"context"

	"github.com/google/uuid"
)

type KeyService interface {
	HandleCreateKey(ctx context.Context, userID uuid.UUID)
}

type keyService struct{}

func NewKeyService() KeyService {
	return &keyService{}
}

func (s *keyService) HandleCreateKey(ctx context.Context, userID uuid.UUID) {
	// bytes := make([]byte, 32)
	// if _, err := rand.Read(bytes); err != nil {

	// }
}
