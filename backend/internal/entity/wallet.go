package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"family_id"`
	Name      string          `gorm:"type:varchar(200);not null" json:"name"`
	Type      string          `gorm:"type:varchar(50);not null" json:"type"`
	Balance   decimal.Decimal `gorm:"type:numeric(18,2);not null;default:0" json:"balance"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
