package models

import "time"

type CLIToken struct {
	BaseModel
	UserID     string    `gorm:"type:uuid;not null;index"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	DeviceName string    `gorm:"not null"`
	LastUsedAt time.Time `gorm:"not null"`
	TokenHash  string    `gorm:"uniqueIndex;not null"`
	ExpiresAt  time.Time `gorm:"not null"`
}
