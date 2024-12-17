package command

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

type UpdateTokenCommand struct {
	Hash     string
	NewToken string
}

type UpdateTokenHandler struct {
	repository mask.Repository
}

func NewUpdateTokenHandler(repository mask.Repository) *UpdateTokenHandler {
	return &UpdateTokenHandler{repository}
}

func (s *UpdateTokenHandler) Handle(ctx context.Context, cmd UpdateTokenCommand) (interface{}, error) {
	return cmd, s.repository.UpdateTokenByHash(ctx, mask.Mask{}, cmd.Hash, cmd.NewToken)
}
