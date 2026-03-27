package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FamilyID uuid.UUID `gorm:"type:uuid;not null;index" json:"family_id"`
	Name     string    `gorm:"type:varchar(200);not null" json:"name"`
	Type     string    `gorm:"type:varchar(50);not null" json:"type"` // income|expense
	Icon     string    `gorm:"type:varchar(100)" json:"icon"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
