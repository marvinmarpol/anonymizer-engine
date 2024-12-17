package query

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

// Query to get cypher by token
type GetRotateCandidateQuery struct {
	Limit         int
	Offset        int
	DayDifference int
	Result        []mask.Mask
}

// Query handler for finding PII by ID
type GetRotateCandidateHandler struct {
	repository mask.Repository
}

func NewGetRotateCandidateHandler(repository mask.Repository) *GetRotateCandidateHandler {
	return &GetRotateCandidateHandler{repository}
}

// Handle executes the query to find PII data by ID
func (h *GetRotateCandidateHandler) Handle(ctx context.Context, q GetRotateCandidateQuery) (GetRotateCandidateQuery, error) {
	mask, err := h.repository.GetRotateCandidate(ctx, q.DayDifference, q.Limit, q.Offset)
	if err != nil {
		return q, err
	}

	q.Result = mask
	return q, err
}
