package router

import (
	"ctx_backend/internal/server/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, h *handler.AuthHandler) {
	r.GET("/auth/:provider", h.Login)
	r.GET("/logout", h.Logout)
	r.GET("/auth/:provider/callback", h.Callback)
	r.POST("/auth/cli/session", h.CreateCLISession)
	r.GET("/auth/cli/session/:session_id", h.PollCLISession)
	r.POST("/auth/refresh", h.RefreshToken)
}
