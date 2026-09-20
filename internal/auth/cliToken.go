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

func (s *CLITokenService) Generate(user models.User, deviceName string) (string, error) {
	token := GenerateRandomString(32)
	rec := models.CLIToken{
		UserID:     user.ID,
		DeviceName: deviceName,
		LastUsedAt: time.Now(),
		TokenHash:  hashToken(token),
		ExpiresAt:  time.Now().Add(90 * 24 * time.Hour),
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
	rec.LastUsedAt = time.Now()
	if err := s.db.Save(&rec).Error; err != nil {
		return "", err
	}
	return rec.UserID, nil
}

func (s *CLITokenService) List(userID string) ([]models.CLIToken, error) {
	var tokens []models.CLIToken
	if err := s.db.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *CLITokenService) Revoke(userID, tokenID string) error {
	result := s.db.Where("user_id = ? AND id = ?", userID, tokenID).Delete(&models.CLIToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}
	return nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
