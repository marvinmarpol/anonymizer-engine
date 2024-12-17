package query

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

// Query to get cypher by token
type GetMaskQuery struct {
	Token  string
	Result mask.Mask
}

// Query handler for finding PII by ID
type GetMaskHandler struct {
	repository mask.Repository
}

func NewGetMaskHandler(repository mask.Repository) *GetMaskHandler {
	return &GetMaskHandler{repository}
}

// Handle executes the query to find PII data by ID
func (h *GetMaskHandler) Handle(ctx context.Context, q GetMaskQuery) (GetMaskQuery, error) {
	mask, err := h.repository.FindByToken(ctx, q.Token)
	if err != nil {
		return q, err
	}

	q.Result = mask
	return q, err
}
