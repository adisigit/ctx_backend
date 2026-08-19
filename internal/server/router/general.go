package router

import (
	"ctx_backend/internal/server/handler"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGeneralRoutes(api huma.API, h *handler.GeneralHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "ctx",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"General"},
	}, h.CtxHandler)

	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/health",
		Tags:        []string{"General"},
	}, h.HealthHandler)
}
