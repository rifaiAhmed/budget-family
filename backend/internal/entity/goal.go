package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Goal struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"family_id"`
	Name          string          `gorm:"type:varchar(200);not null" json:"name"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"target_amount"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(18,2);not null;default:0" json:"current_amount"`
	TargetDate    time.Time       `gorm:"type:date" json:"target_date"`
}

func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
