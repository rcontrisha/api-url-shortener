package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (l *LinkRepo) SetShortenLink(ctx context.Context, userId int, originalUrl, slug string) error {
	query := `insert into links (user_id, original_url, slug) values ($1, $2, $3)`

	_, err := l.db.Exec(ctx, query, userId, originalUrl, slug)
	if err != nil {
		var pgErr *pgconn.PgError
		// Check if the error is a Postgres error and has code "23505" (Duplicate Entry)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("slug '%s' already exists (duplicate key)", slug)
		}
		return err
	}
	return nil
}

func (l *LinkRepo) GetOriginalUrl(ctx context.Context, slug string) (string, error) {
	query := `select original_url from links where slug = $1 and deleted_at is null`

	var originalUrl string
	if err := l.db.QueryRow(ctx, query, slug).Scan(&originalUrl); err != nil {
		return "", err
	}

	return originalUrl, nil
}

func (l *LinkRepo) DeleteShortenedLink(ctx context.Context, id int, userId int) error {
	query := `
		UPDATE links 
		SET deleted_at = NOW() 
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	_, err := l.db.Exec(ctx, query, id, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("link not found or unauthorized")
		}
		return err
	}

	return nil
}
