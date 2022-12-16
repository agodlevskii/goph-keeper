package user

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
	CreateUserTable = `CREATE TABLE IF NOT EXISTS users(
    	id UUID DEFAULT gen_random_uuid(),
    	name VARCHAR(255),
    	password VARCHAR(255),
    	UNIQUE(name),
    	PRIMARY KEY(id))`
	AddUser       = "INSERT INTO users(name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id"
	DeleteUser    = "DELETE FROM users WHERE id = $1"
	GetUserByID   = "SELECT * FROM users WHERE id = $1"
	GetUserByName = "SELECT * FROM users WHERE name = $1"
)

func NewDBRepo(url string) (*DBRepo, error) {
	if url == "" {
		return &DBRepo{}, ErrDBMissingURL
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &DBRepo{}, err
	}

	_, err = db.ExecContext(context.Background(), CreateUserTable)
	return &DBRepo{db: db}, err
}

func (r *DBRepo) AddUser(ctx context.Context, user User) (User, error) {
	var id string
	if err := r.db.QueryRowContext(ctx, AddUser, user.Name, user.Password).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	user.ID = id
	return user, nil
}

func (r *DBRepo) DeleteUser(ctx context.Context, uid string) error {
	_, err := r.db.ExecContext(ctx, DeleteUser, uid)
	return err
}

func (r *DBRepo) GetUserByID(ctx context.Context, uid string) (User, error) {
	return r.getUser(ctx, GetUserByID, uid)
}

func (r *DBRepo) GetUserByName(ctx context.Context, name string) (User, error) {
	return r.getUser(ctx, GetUserByName, name)
}

func (r *DBRepo) getUser(ctx context.Context, query string, args ...any) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, nil
}
