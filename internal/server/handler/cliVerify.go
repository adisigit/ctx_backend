package handler

import (
	"context"
	"ctx_backend/internal/auth"
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

type CLITokenInfo struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
	CreatedAt  string `json:"created_at"`
	LastUsedAt string `json:"last_used_at"`
	ExpiresAt  string `json:"expires_at"`
}

type ListCLITokensOutput struct {
	Body struct {
		Tokens []CLITokenInfo `json:"tokens"`
	}
}

func ListCLITokens(cliTokenService *auth.CLITokenService) func(context.Context, *struct{}) (*ListCLITokensOutput, error) {
	return func(ctx context.Context, input *struct{}) (*ListCLITokensOutput, error) {
		userID, _ := ctx.Value(middleware.ContextKeyUserID).(string)
		token, err := cliTokenService.List(userID)
		if err != nil {
			return nil, err
		}
		resp := &ListCLITokensOutput{}
		for _, t := range token {
			resp.Body.Tokens = append(resp.Body.Tokens, CLITokenInfo{
				ID:         t.ID,
				DeviceName: t.DeviceName,
				CreatedAt:  t.CreatedAt.String(),
				LastUsedAt: t.LastUsedAt.String(),
				ExpiresAt:  t.ExpiresAt.String(),
			})
		}
		return resp, nil
	}
}

type RevokeCLITokenInput struct {
	TokenID string `path:"id"`
}

type RevokeCLITokenOutput struct {
	Body struct {
		Success bool `json:"success"`
	}
}

func RevokeCLIToken(cliTokenService *auth.CLITokenService) func(context.Context, *RevokeCLITokenInput) (*RevokeCLITokenOutput, error) {
	return func(ctx context.Context, input *RevokeCLITokenInput) (*RevokeCLITokenOutput, error) {
		userID, _ := ctx.Value(middleware.ContextKeyUserID).(string)
		err := cliTokenService.Revoke(userID, input.TokenID)
		if err != nil {
			return nil, err
		}
		resp := &RevokeCLITokenOutput{}
		resp.Body.Success = true
		return resp, nil
	}
}
