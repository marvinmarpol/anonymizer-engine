package command

import "context"

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) (interface{}, error)
}

type Commands struct {
	CreateMaskCommand  CommandHandler[CreateMaskCommand]
	UpdateTokenCommand CommandHandler[UpdateTokenCommand]
	UpdateMaskCommand  CommandHandler[UpdateMaskCommand]
}
