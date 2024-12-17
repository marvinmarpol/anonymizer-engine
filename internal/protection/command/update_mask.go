package command

import (
	"context"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
)

type UpdateMaskCommand mask.Mask

type UpdateMaskHandler struct {
	repository mask.Repository
}

func NewUpdateMaskHandler(repository mask.Repository) *UpdateMaskHandler {
	return &UpdateMaskHandler{repository}
}

func (s *UpdateMaskHandler) Handle(ctx context.Context, cmd UpdateMaskCommand) (interface{}, error) {
	return cmd, s.repository.Update(ctx, mask.Mask(cmd))
}
