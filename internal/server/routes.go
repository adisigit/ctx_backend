package server

import (
	"ctx_backend/internal/auth"
	"ctx_backend/internal/server/handler"
	"ctx_backend/internal/server/router"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	api := humagin.New(r, buildHumaConfig())
	cliTokenService := auth.NewCLITokenService(s.db.DB())
	cliSessionService := auth.NewCLISessionService()
	jwtService := auth.NewJWTService()
	router.RegisterGeneralRoutes(api, handler.NewGeneralHandler(s.db))
	router.RegisterAuthRoutes(r, handler.NewAuthHandler(s.db.DB(), cliTokenService, cliSessionService, jwtService))
	router.RegisterProtectedGroup(api, jwtService, cliTokenService)
	return r
}

func buildHumaConfig() huma.Config {
	config := huma.DefaultConfig("CTX API", "1.0.0")
	config.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"bearer": {
				Type: "http", Scheme: "bearer", BearerFormat: "JWT",
			},
		},
	}
	return config
}
