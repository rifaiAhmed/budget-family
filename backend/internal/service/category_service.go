package service

import (
	"context"
	"strings"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CategoryService interface {
	Create(ctx context.Context, familyID uuid.UUID, name, typ, icon string) (*entity.Category, error)
	List(ctx context.Context, familyID uuid.UUID, typ string, page, limit int) ([]entity.Category, utils.PageMeta, error)
}

type categoryService struct {
	cfg    config.Config
	logger *zap.Logger
	repo   repository.CategoryRepository
}

func NewCategoryService(cfg config.Config, logger *zap.Logger, repo repository.CategoryRepository) CategoryService {
	return &categoryService{cfg: cfg, logger: logger, repo: repo}
}

func (s *categoryService) Create(ctx context.Context, familyID uuid.UUID, name, typ, icon string) (*entity.Category, error) {
	c := &entity.Category{ID: uuid.New(), FamilyID: familyID, Name: strings.TrimSpace(name), Type: strings.TrimSpace(typ), Icon: strings.TrimSpace(icon)}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, utils.NewInternal("failed to create category", err)
	}
	return c, nil
}

func (s *categoryService) List(ctx context.Context, familyID uuid.UUID, typ string, page, limit int) ([]entity.Category, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	items, total, err := s.repo.ListByFamily(ctx, familyID, typ, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list categories", err)
	}
	return items, utils.BuildMeta(page, limit, total), nil
}
