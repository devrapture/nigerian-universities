package service

import (
	"context"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
)

type institutionService struct {
	repo repositories.InstitutionRepository
}

type InstitutionService interface {
	StoreScrapedInstitutions(ctx context.Context, institutions []model.Institution) error
}

func NewInstitutionService(repo repositories.InstitutionRepository) InstitutionService {
	return &institutionService{
		repo: repo,
	}
}

func (s *institutionService) StoreScrapedInstitutions(ctx context.Context, institutions []model.Institution) error {
	return s.repo.UpsertMany(ctx, institutions)
}
