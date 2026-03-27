package repository

import (
	"context"

	"budget-family/internal/entity"

	"gorm.io/gorm"
)

type InvitationRepository interface {
	Create(ctx context.Context, inv *entity.Invitation) error
}

type invitationRepository struct{ db *gorm.DB }

func NewInvitationRepository(db *gorm.DB) InvitationRepository { return &invitationRepository{db: db} }

func (r *invitationRepository) Create(ctx context.Context, inv *entity.Invitation) error {
	return r.db.WithContext(ctx).Create(inv).Error
}
