package handler

import (
	"ctx_backend/internal/auth"
	"ctx_backend/internal/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Login(c *gin.Context) {
	providerName := c.Param("provider")
	provider, err := auth.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider not supported"})
		return
	}
	url := provider.Config().AuthCodeURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	providerName := c.Param("provider")
	provider, err := auth.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider not supported"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not provided"})
		return
	}
	token, err := provider.Config().Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to exchange code for token"})
		return
	}
	providerUser, err := provider.GetUser(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get user"})
		return
	}
	var user models.User
	if err := h.db.Where(models.User{ProviderID: providerUser.ID, Provider: providerName}).Attrs(models.User{
		Name:   providerUser.Name,
		Email:  providerUser.Email,
		Avatar: providerUser.Avatar,
	}).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create or find user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
}
