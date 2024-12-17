package service

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/entity"
)

type Services interface {
	Deidentify(ctx context.Context, cmd interface{}) (interface{}, error)
	Reidentify(ctx context.Context, cmd interface{}) (interface{}, error)
	GetCypher(ctx context.Context, cmd entity.GetCypherPayload) (interface{}, error)
	RotateKeys(ctx context.Context, cmd entity.RotatePayload) (interface{}, error)
}
