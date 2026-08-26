package handler

import (
	"crypto/subtle"
	"ctx_backend/internal/auth"
	"ctx_backend/internal/database/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db                *gorm.DB
	cliTokenService   *auth.CLITokenService
	cliSessionService *auth.CLISessionService
	jwtService        *auth.JWTService
	frontendURL       string
}

func NewAuthHandler(db *gorm.DB, cliTokenService *auth.CLITokenService, cliSessionService *auth.CLISessionService, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		db:                db,
		cliTokenService:   cliTokenService,
		cliSessionService: cliSessionService,
		jwtService:        jwtService,
		frontendURL:       os.Getenv("FRONTEND_URL"),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	providerName := c.Param("provider")
	provider, err := auth.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider not supported"})
		return
	}
	client := c.DefaultQuery("client", "web")
	sessionID := c.Query("session_id")
	if client == "cli" && sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required for cli client"})
		return
	}
	csrf := auth.GenerateRandomString(32)
	c.SetCookie("csrf", csrf, 300, "/", "", isProduction(), true)

	state, err := auth.EncodeState(auth.OAuthState{Client: client, SessionID: sessionID, CSRF: csrf})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode state"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, provider.Config().AuthCodeURL(state))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", isProduction(), true)
	c.SetCookie("refresh_token", "", -1, "/", "", isProduction(), true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	providerName := c.Param("provider")
	provider, err := auth.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider not supported"})
		return
	}
	state, err := auth.DecodeState(c.Query("state"))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	csrfCookie, err := c.Cookie("csrf")
	if err != nil || csrfCookie == "" ||
		subtle.ConstantTimeCompare([]byte(csrfCookie), []byte(state.CSRF)) != 1 {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	c.SetCookie("csrf", "", -1, "/", "", isProduction(), true)
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	token, err := provider.Config().Exchange(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	providerUser, err := provider.GetUser(token)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	var user models.User
	if err := h.db.Where(models.User{ProviderID: providerUser.ID, Provider: providerName}).Attrs(models.User{
		Name:   providerUser.Name,
		Email:  providerUser.Email,
		Avatar: providerUser.Avatar,
	}).FirstOrCreate(&user).Error; err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?status=error")
		return
	}
	switch state.Client {
	case "cli":
		h.handleCLICallback(c, user, state.SessionID)
	default:
		h.handleWebCallback(c, user)
	}
}

func (h *AuthHandler) handleCLICallback(c *gin.Context, user models.User, sessionID string) {
	accessToken, err := h.cliTokenService.Generate(user)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=cli&status=error")
		return
	}
	if err := h.cliSessionService.Complete(sessionID, accessToken); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=cli&status=error")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=cli&status=success")
}

func (h *AuthHandler) handleWebCallback(c *gin.Context, user models.User) {
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Email, user.Name)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=web&status=error")
		return
	}
	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Email, user.Name)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=web&status=error")
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 15*60, "/", "", isProduction(), true)
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", isProduction(), true)
	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login/complete?client=web&status=success")
}

func (h *AuthHandler) CreateCLISession(c *gin.Context) {
	var req struct {
		PublicKey string `json:"public_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessionID := h.cliSessionService.Create(req.PublicKey)
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"login_url":  h.frontendURL + "/login?client=cli&session_id=" + sessionID,
	})
}

func (h *AuthHandler) PollCLISession(c *gin.Context) {
	sess, err := h.cliSessionService.Poll(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sess.Status != auth.CLISessionCompleted {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          string(sess.Status),
		"encrypted_token": sess.EncryptedToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil || refreshCookie == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token not found"})
		return
	}
	userId, err := h.jwtService.Verify(refreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	var user models.User
	if err := h.db.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newAccessToken, err := h.jwtService.GenerateAccessToken(userId, user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("access_token", newAccessToken, 15*60, "/", "", isProduction(), true)
	c.JSON(http.StatusOK, gin.H{"status": "refreshed"})
}

func isProduction() bool {
	return os.Getenv("MODE") == "production"
}
