package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
