package db

import (
	"Avito/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context, config config.StorageConfig) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn(config))
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn(config config.StorageConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)
}
