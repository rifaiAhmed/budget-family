package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Budget struct {
	ID         uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID   uuid.UUID       `gorm:"type:uuid;not null;index;uniqueIndex:uniq_budget" json:"family_id"`
	CategoryID uuid.UUID       `gorm:"type:uuid;not null;index;uniqueIndex:uniq_budget" json:"category_id"`
	Amount     decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	Month      int             `gorm:"not null;uniqueIndex:uniq_budget" json:"month"`
	Year       int             `gorm:"not null;uniqueIndex:uniq_budget" json:"year"`
}

func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
