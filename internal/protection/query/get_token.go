package query

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

// Query to get cypher by token
type GetTokenQuery struct {
	Hash   string
	Result string
}

// Query handler for finding PII by ID
type GetTokenHandler struct {
	repository mask.Repository
}

func NewGetTokenHandler(repository mask.Repository) *GetTokenHandler {
	return &GetTokenHandler{repository}
}

// Handle executes the query to find PII data by ID
func (h *GetTokenHandler) Handle(ctx context.Context, q GetTokenQuery) (GetTokenQuery, error) {
	token, err := h.repository.GetTokenByHash(ctx, q.Hash)
	if err != nil {
		return q, err
	}

	q.Result = token
	return q, err
}
