package mask

import (
	"context"

	"github.com/go-pg/pg/v10"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *pg.DB
}

// NewPostgresRepository creates a new repository instance
func NewPostgresRepository(db *pg.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create adds a new entity to the database
func (r *PostgresRepository) Create(ctx context.Context, entity Mask) error {
	_, err := r.db.Model(&entity).Insert()
	return err
}

// Update the mask entity
func (r *PostgresRepository) Update(ctx context.Context, entity Mask) error {
	_, err := r.db.Model(&entity).
		Set("key = ?", entity.Key).
		Set("cypher = ?", entity.Cypher).
		Set("rotated_at = ?", "now()").
		Where("token = ?", entity.Token).
		Update()

	return err
}

// UpdateTokenByHash modifies existing token by old token
func (r *PostgresRepository) UpdateToken(ctx context.Context, entity Mask, oldToken, newToken string) error {
	_, err := r.db.Model(&entity).
		Set("token = ?", newToken).
		Set("updated_at = ?", "now()").
		Where("token = ?", oldToken).
		Update()

	return err
}

// UpdateTokenByHash modifies existing token  by hash
func (r *PostgresRepository) UpdateTokenByHash(ctx context.Context, entity Mask, hash, newToken string) error {
	_, err := r.db.Model(&entity).
		Set("token = ?", newToken).
		Set("updated_at = ?", "now()").
		Where("hash = ?", hash).
		Update()

	return err
}

// FindByToken retrieves an entity by token
func (r *PostgresRepository) FindByToken(ctx context.Context, token string) (Mask, error) {
	var entity Mask

	err := r.db.ModelContext(ctx, &entity).
		Where("token = ?", token).
		Select()

	return entity, err
}

// GetTokenByHash retrieves unique token by hash
func (r *PostgresRepository) GetTokenByHash(ctx context.Context, hash string) (string, error) {
	var entity Mask

	err := r.db.ModelContext(ctx, &entity).
		Column("token").
		Where("hash = ?", hash).
		Select()

	return entity.Token, err
}

// GetTokenByHash retrieves unique token by hash
func (r *PostgresRepository) GetRotateCandidate(ctx context.Context, dayDifference, limit, offset int) ([]Mask, error) {
	var entities []Mask

	err := r.db.ModelContext(ctx, &entities).
		Where("rotated_at < NOW() - INTERVAL '? days'", dayDifference).
		Order("rotated_at DESC").
		Limit(limit).
		Offset(offset).
		Select()

	return entities, err
}
