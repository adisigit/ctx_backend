package router

import (
	"ctx_backend/internal/auth"
	"ctx_backend/internal/server/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterProtectedGroup(api huma.API, jwtService *auth.JWTService, cliTokenService *auth.CLITokenService) {
	group := huma.NewGroup(api, "/api")
	group.UseMiddleware(middleware.RequireAuth(jwtService, cliTokenService))
	registerCLIVerify(group)
}
