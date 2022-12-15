package session

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver
)

type DBRepo struct {
	db *sql.DB
}

const (
	CreateSessionTable = `CREATE TABLE IF NOT EXISTS sessions(
    	cid VARCHAR(50),
	   	token VARCHAR(165),
	   	PRIMARY KEY (cid)
	)`
	DeleteSession = `DELETE FROM sessions WHERE cid = $1`
	GetSession    = "SELECT token FROM sessions WHERE cid = $1"
	StoreSession  = "INSERT INTO sessions(cid, token) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING token"
)

func NewDBRepo(url string) (*DBRepo, error) {
	if url == "" {
		return &DBRepo{}, ErrDBMissingURL
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &DBRepo{}, err
	}

	_, err = db.ExecContext(context.Background(), CreateSessionTable)
	return &DBRepo{db: db}, nil
}

func (r *DBRepo) DeleteSession(ctx context.Context, cid string) error {
	_, err := r.db.ExecContext(ctx, DeleteSession, cid)
	return err
}

func (r *DBRepo) GetSession(ctx context.Context, cid string) (string, error) {
	var token string
	err := r.db.QueryRowContext(ctx, GetSession, cid).Scan(&token)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return token, err
}

func (r *DBRepo) StoreSession(ctx context.Context, cid, token string) error {
	_, err := r.db.ExecContext(ctx, StoreSession, cid, token)
	return err
}
