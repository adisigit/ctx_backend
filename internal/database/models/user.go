package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name       string `gorm:"not null"`
	Email      string `gorm:"not null;unique"`
	Avatar     string
	Provider   string `gorm:"not null"`
	ProviderID string `gorm:"uniqueIndex;not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil.String() {
		u.ID = uuid.New().String()
	}
	return nil
}
