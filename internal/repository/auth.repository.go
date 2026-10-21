package repository

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-url-shortener/internal/model"
	"github.com/redis/go-redis/v9"
)

type AuthRepo struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepo(db *pgxpool.Pool, rdb *redis.Client) *AuthRepo {
	return &AuthRepo{
		db:  db,
		rdb: rdb,
	}
}

func (a *AuthRepo) FindAccount(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash
		FROM users u
		WHERE u.email = $1`

	var acc model.User
	err := a.db.QueryRow(ctx, query, email).Scan(&acc.Id, &acc.Email, &acc.Password)
	if err != nil {
		return model.User{}, err
	}
	return acc, nil
}

func (a *AuthRepo) CreateAccount(ctx context.Context, payload model.User) error {
	query := `insert into users (email, password_hash) values ($1, $2)`

	cmd, err := a.db.Exec(ctx, query, payload.Email, payload.Password)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

func (a *AuthRepo) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := os.Getenv("REDIS_USER") + ":blacklist:" + jti
	return a.rdb.Set(ctx, key, "revoked", ttl).Err()
}
