package router

import (
	"ctx_backend/internal/server/handler"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(api huma.API, r *gin.Engine, h *handler.AuthHandler) {
	r.GET("/auth/:provider", h.Login)
	r.GET("/logout", h.Logout)
	r.GET("/auth/:provider/callback", h.Callback)
}
