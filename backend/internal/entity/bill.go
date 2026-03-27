package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Bill struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"family_id"`
	Name      string          `gorm:"type:varchar(200);not null" json:"name"`
	Amount    decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	DueDay    int             `gorm:"not null" json:"due_day"`
	Recurring bool            `gorm:"not null;default:false" json:"recurring"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type Invitation struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID uuid.UUID `gorm:"type:uuid;not null;index" json:"family_id"`
	Email    string    `gorm:"type:varchar(255);not null;index" json:"email"`
	Token    string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	Status   string    `gorm:"type:varchar(50);not null" json:"status"`
}

func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
