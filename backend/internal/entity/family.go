package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Family struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (f *Family) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type FamilyMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:uniq_family_user" json:"family_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:uniq_family_user" json:"user_id"`
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (m *FamilyMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
