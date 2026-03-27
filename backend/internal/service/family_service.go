package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FamilyService interface {
	Create(ctx context.Context, ownerID uuid.UUID, name string) (*entity.Family, error)
	List(ctx context.Context, userID uuid.UUID) ([]entity.Family, error)
	Join(ctx context.Context, userID, familyID uuid.UUID) (*entity.Family, error)
	Invite(ctx context.Context, requesterID uuid.UUID, familyID uuid.UUID, email string) (*entity.Invitation, error)
	Members(ctx context.Context, requesterID, familyID uuid.UUID) ([]repository.FamilyMemberRow, error)
}

type familyService struct {
	cfg         config.Config
	logger      *zap.Logger
	families    repository.FamilyRepository
	invitations repository.InvitationRepository
}

func NewFamilyService(cfg config.Config, logger *zap.Logger, families repository.FamilyRepository, invitations repository.InvitationRepository) FamilyService {
	return &familyService{cfg: cfg, logger: logger, families: families, invitations: invitations}
}

func (s *familyService) Create(ctx context.Context, ownerID uuid.UUID, name string) (*entity.Family, error) {
	f := &entity.Family{ID: uuid.New(), Name: strings.TrimSpace(name), OwnerID: ownerID}
	owner := &entity.FamilyMember{ID: uuid.New(), FamilyID: f.ID, UserID: ownerID, Role: "owner"}
	if err := s.families.CreateFamily(ctx, f, owner); err != nil {
		return nil, utils.NewInternal("failed to create family", err)
	}
	return f, nil
}

func (s *familyService) List(ctx context.Context, userID uuid.UUID) ([]entity.Family, error) {
	families, err := s.families.ListFamiliesByUser(ctx, userID)
	if err != nil {
		return nil, utils.NewInternal("failed to list families", err)
	}
	return families, nil
}

func (s *familyService) Join(ctx context.Context, userID, familyID uuid.UUID) (*entity.Family, error) {
	f, err := s.families.GetByID(ctx, familyID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "record not found") {
			return nil, utils.NewNotFound("family not found", err)
		}
		return nil, utils.NewInternal("failed to get family", err)
	}

	isMember, err := s.families.IsMember(ctx, familyID, userID)
	if err != nil {
		return nil, utils.NewInternal("failed to verify membership", err)
	}
	if isMember {
		return f, nil
	}

	m := &entity.FamilyMember{ID: uuid.New(), FamilyID: familyID, UserID: userID, Role: "member"}
	if err := s.families.AddMember(ctx, m); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return f, nil
		}
		return nil, utils.NewInternal("failed to join family", err)
	}
	return f, nil
}

func (s *familyService) Invite(ctx context.Context, requesterID uuid.UUID, familyID uuid.UUID, email string) (*entity.Invitation, error) {
	isMember, err := s.families.IsMember(ctx, familyID, requesterID)
	if err != nil {
		return nil, utils.NewInternal("failed to verify membership", err)
	}
	if !isMember {
		return nil, utils.NewForbidden("not a family member", nil)
	}

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, utils.NewBadRequest("email is required", nil)
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, utils.NewInternal("failed to generate invitation token", err)
	}
	token := hex.EncodeToString(b)

	inv := &entity.Invitation{ID: uuid.New(), FamilyID: familyID, Email: email, Token: token, Status: "pending"}
	if err := s.invitations.Create(ctx, inv); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, utils.NewConflict("invitation already exists", err)
		}
		return nil, utils.NewInternal("failed to create invitation", err)
	}
	return inv, nil
}

func (s *familyService) Members(ctx context.Context, requesterID, familyID uuid.UUID) ([]repository.FamilyMemberRow, error) {
	isMember, err := s.families.IsMember(ctx, familyID, requesterID)
	if err != nil {
		return nil, utils.NewInternal("failed to verify membership", err)
	}
	if !isMember {
		return nil, utils.NewForbidden("not a family member", nil)
	}

	rows, err := s.families.ListMembers(ctx, familyID)
	if err != nil {
		return nil, utils.NewInternal("failed to list members", err)
	}
	return rows, nil
}
