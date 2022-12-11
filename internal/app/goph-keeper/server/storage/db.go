package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver
	log "github.com/sirupsen/logrus"
)

const (
	CreateUserTable = `CREATE TABLE IF NOT EXISTS users(
    	id UUID DEFAULT gen_random_uuid(),
    	name VARCHAR(255),
    	password VARCHAR(255),
    	UNIQUE(name),
    	PRIMARY KEY(id))`
	CreateStorageTable = `CREATE TABLE IF NOT EXISTS storage(
    	id UUID DEFAULT gen_random_uuid(),
    	uid UUID,
    	data BYTEA,
    	type INT,
    	PRIMARY KEY(id),
		CONSTRAINT fk_user
		    FOREIGN KEY (uid)
		        REFERENCES users(id))`
	CreateSessionTable = `CREATE TABLE IF NOT EXISTS sessions(
    	cid VARCHAR(50),
	   	token VARCHAR(165),
	   	PRIMARY KEY (cid)
	)`
	DeleteSession    = `DELETE FROM sessions WHERE cid = $1`
	GetSession       = "SELECT token FROM sessions WHERE cid = $1"
	StoreSession     = "INSERT INTO sessions(cid, token) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING token"
	AddUser          = "INSERT INTO users(name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id"
	GetUserByID      = "SELECT * FROM users WHERE id = $1"
	GetUserByName    = "SELECT * FROM users WHERE name = $1"
	GetAllDataByType = "SELECT * FROM storage WHERE uid = $1 AND type = $2"
	GetDataByID      = "SELECT * FROM storage WHERE uid = $1 AND id = $2"
	StoreData        = `INSERT INTO storage(uid, data, type) VALUES($1, $2, $3)
																 ON CONFLICT DO NOTHING RETURNING id`
)

type DBRepo struct {
	db *sql.DB
}

func NewDBRepo(ctx context.Context, url string) (*DBRepo, error) {
	if url == "" {
		return &DBRepo{}, errors.New("db url is missing")
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &DBRepo{}, err
	}

	if err = initDB(ctx, db); err != nil {
		return &DBRepo{}, err
	}
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
		err = errors.New("token not found")
	}
	return token, err
}

func (r *DBRepo) StoreSession(ctx context.Context, cid, token string) error {
	_, err := r.db.ExecContext(ctx, StoreSession, cid, token)
	return err
}

func (r *DBRepo) AddUser(ctx context.Context, u User) (User, error) {
	var id string
	if err := r.db.QueryRowContext(ctx, AddUser, u.Name, u.Password).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("user with the specified name already exists")
		}
		return User{}, err
	}

	u.ID = id
	return u, nil
}

func (r *DBRepo) GetUserByID(ctx context.Context, uid string) (User, error) {
	return r.getUser(ctx, GetUserByID, uid)
}

func (r *DBRepo) GetUserByName(ctx context.Context, name string) (User, error) {
	return r.getUser(ctx, GetUserByName, name)
}

func (r *DBRepo) GetAllDataByType(ctx context.Context, uid string, t Type) ([]SecureData, error) {
	rows, err := r.db.QueryContext(ctx, GetAllDataByType, uid, t)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer closeRows(rows)

	var data []SecureData
	for rows.Next() {
		var piece SecureData
		if err = rows.Scan(&piece.ID, &piece.UID, &piece.Data, &piece.Type); err != nil {
			return nil, err
		}
		data = append(data, piece)
	}
	return data, nil
}

func (r *DBRepo) GetDataByID(ctx context.Context, uid, id string) (SecureData, error) {
	var data SecureData
	err := r.db.QueryRowContext(ctx, GetDataByID, uid, id).Scan(&data.ID, &data.UID, &data.Data, &data.Type)
	return data, err
}

func (r *DBRepo) StoreData(ctx context.Context, data SecureData) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, StoreData, data.UID, data.Data, data.Type).Scan(&id)
	return id, err
}

func (r *DBRepo) getUser(ctx context.Context, query string, args ...any) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return user, nil
}

func initDB(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, CreateUserTable); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, CreateStorageTable); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, CreateSessionTable); err != nil {
		return err
	}
	return nil
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Error(err)
	}
}
