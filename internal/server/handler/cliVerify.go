package handler

import (
	"context"
	"ctx_backend/internal/server/middleware"
)

type VerifyCLITokenOutput struct {
	Body struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"user_id"`
	}
}

func VerifyCLIToken(ctx context.Context, input *struct{}) (*VerifyCLITokenOutput, error) {
	userID, _ := ctx.Value(middleware.ContextKeyUserID).(string)
	resp := &VerifyCLITokenOutput{}
	resp.Body.Valid = userID != ""
	resp.Body.UserID = userID
	return resp, nil
}
