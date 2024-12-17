package mask

import "context"

// mockgen -source=internal/protection/domain/mask/repository.go -destination=internal/protection/domain/mask/mock/repository.go
// Repository interface defines the methods to interact with the database
type Repository interface {
	Create(ctx context.Context, entity Mask) error
	Update(ctx context.Context, entity Mask) error
	UpdateToken(ctx context.Context, entity Mask, oldToken, newToken string) error
	UpdateTokenByHash(ctx context.Context, entity Mask, hash, newToken string) error
	FindByToken(ctx context.Context, token string) (Mask, error)
	GetTokenByHash(ctx context.Context, hash string) (string, error)
	GetRotateCandidate(ctx context.Context, dayDifference, limit, offset int) ([]Mask, error)
}
