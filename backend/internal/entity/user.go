package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(200);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(30)" json:"phone"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
