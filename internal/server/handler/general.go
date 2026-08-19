package handler

import (
	"context"
	"ctx_backend/internal/database"
)

type GeneralHandler struct {
	db database.Service
}

func NewGeneralHandler(db database.Service) *GeneralHandler {
	return &GeneralHandler{db: db}
}

type CtxOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (h *GeneralHandler) CtxHandler(ctx context.Context, input *struct{}) (*CtxOutput, error) {
	resp := &CtxOutput{}
	resp.Body.Message = "CTX"
	return resp, nil
}

type HealthOutput struct {
	Body map[string]string `json:"body"`
}

func (h *GeneralHandler) HealthHandler(ctx context.Context, input *struct{}) (*HealthOutput, error) {
	resp := &HealthOutput{}
	resp.Body = h.db.Health()
	return resp, nil
}
