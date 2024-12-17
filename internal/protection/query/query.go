package query

import (
	"context"
)

// Generic Query Handler Interface
type QueryHandler[Q any] interface {
	Handle(ctx context.Context, q Q) (Q, error)
}

type Queries struct {
	GetCypherQuery          QueryHandler[GetCypherQuery]
	GetMaskQuery            QueryHandler[GetMaskQuery]
	GetTokenQuery           QueryHandler[GetTokenQuery]
	GetRotateCandidateQuery QueryHandler[GetRotateCandidateQuery]
}
