package router

import (
	"ctx_backend/internal/auth"
	"ctx_backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerCLIVerify(group huma.API, cliTokenService *auth.CLITokenService) {
	huma.Register(group, huma.Operation{
		OperationID: "verify-cli-token",
		Method:      http.MethodGet,
		Path:        "/auth/cli/verify",
		Tags:        []string{"Auth"},
		Summary:     "check token/session still valid",
	}, handler.VerifyCLIToken)

	huma.Register(group, huma.Operation{
		OperationID: "list-cli-token",
		Method:      http.MethodGet,
		Path:        "/auth/cli/list",
		Tags:        []string{"Auth"},
		Summary:     "list cli token",
	}, handler.ListCLITokens(cliTokenService))

	huma.Register(group, huma.Operation{
		OperationID: "revoke-cli-token",
		Method:      http.MethodDelete,
		Path:        "/auth/cli/revoke/{id}",
		Tags:        []string{"Auth"},
		Summary:     "revoke a cli token",
	}, handler.RevokeCLIToken(cliTokenService))
}
