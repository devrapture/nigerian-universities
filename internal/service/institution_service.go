package service

import (
	"context"

	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
)

type institutionService struct {
	repo repositories.InstitutionRepository
}

type InstitutionService interface {
	StoreScrapedInstitutions(ctx context.Context, institutions []model.Institution) error
	GetAllInstitutions(ctx context.Context, queryDTO dto.ListInstitutionQuery) ([]model.Institution, int64, error)
}

func NewInstitutionService(repo repositories.InstitutionRepository) InstitutionService {
	return &institutionService{
		repo: repo,
	}
}

func (s *institutionService) StoreScrapedInstitutions(ctx context.Context, institutions []model.Institution) error {
	return s.repo.UpsertMany(ctx, institutions)
}

func (s *institutionService) GetAllInstitutions(ctx context.Context, queryDTO dto.ListInstitutionQuery) ([]model.Institution, int64, error) {
	return s.repo.FindAll(ctx, queryDTO)
}
