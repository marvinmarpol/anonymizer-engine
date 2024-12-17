package command

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

type CreateMaskCommand mask.Mask
type CreateMaskHandler struct {
	repository mask.Repository
}

func NewCreateMaskHandler(repository mask.Repository) *CreateMaskHandler {
	return &CreateMaskHandler{repository}
}

func (s *CreateMaskHandler) Handle(ctx context.Context, cmd CreateMaskCommand) (interface{}, error) {
	return cmd, s.repository.Create(ctx, mask.Mask(cmd))
}
