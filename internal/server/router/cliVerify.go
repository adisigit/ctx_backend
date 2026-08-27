package router

import (
	"ctx_backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func registerCLIVerify(group huma.API) {
	huma.Register(group, huma.Operation{
		OperationID: "verify-cli-token",
		Method:      http.MethodGet,
		Path:        "/auth/cli/verify",
		Tags:        []string{"Auth"},
		Summary:     "check token/session still valid",
	}, handler.VerifyCLIToken)
}
