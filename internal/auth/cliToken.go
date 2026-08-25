package auth

import (
	"crypto/sha256"
	"ctx_backend/internal/database/models"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CLITokenService struct {
	db *gorm.DB
}

func NewCLITokenService(db *gorm.DB) *CLITokenService {
	return &CLITokenService{db: db}
}

func (s *CLITokenService) Generate(user models.User) (string, error) {
	token := GenerateRandomString(32)
	rec := models.CLIToken{
		UserID:    user.ID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
	}
	if err := s.db.Create(&rec).Error; err != nil {
		return "", err
	}
	return token, nil
}

func (s *CLITokenService) Verify(token string) (string, error) {
	var rec models.CLIToken
	if err := s.db.Where("token_hash = ?", hashToken(token)).First(&rec).Error; err != nil {
		return "", err
	}
	if time.Now().After(rec.ExpiresAt) {
		return "", errors.New("token expired")
	}
	return rec.UserID, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
