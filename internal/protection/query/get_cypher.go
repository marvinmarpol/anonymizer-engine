package query

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

// Query to get cypher by token
type GetCypherQuery struct {
	Token  string
	Result string
}

// Query handler for finding PII by ID
type GetCypherHandler struct {
	repository mask.Repository
}

func NewGetCypherHandler(repository mask.Repository) *GetCypherHandler {
	return &GetCypherHandler{repository}
}

// Handle executes the query to find PII data by ID
func (h *GetCypherHandler) Handle(ctx context.Context, q GetCypherQuery) (GetCypherQuery, error) {
	mask, err := h.repository.FindByToken(ctx, q.Token)
	if err != nil {
		return q, err
	}

	q.Result = mask.Cypher
	return q, err
}
