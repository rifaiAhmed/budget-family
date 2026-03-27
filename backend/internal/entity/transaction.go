package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"family_id"`
	WalletID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"wallet_id"`
	CategoryID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"category_id"`
	Amount          decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	Type            string          `gorm:"type:varchar(50);not null" json:"type"` // income|expense
	Note            string          `gorm:"type:text" json:"note"`
	TransactionDate time.Time       `gorm:"type:date;not null;index" json:"transaction_date"`
	CreatedBy       uuid.UUID       `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedByName   string          `gorm:"-" json:"created_by_name"`
	CreatedByEmail  string          `gorm:"-" json:"created_by_email"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
