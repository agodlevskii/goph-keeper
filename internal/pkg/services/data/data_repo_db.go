package data

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver
	log "github.com/sirupsen/logrus"
)

type DBRepo struct {
	db *sql.DB
}

const (
	CreateStorageTable = `CREATE TABLE IF NOT EXISTS storage(
    	id UUID DEFAULT gen_random_uuid(),
    	uid UUID,
    	data BYTEA,
    	type INT,
    	PRIMARY KEY(id),
		CONSTRAINT fk_user
		    FOREIGN KEY (uid)
		        REFERENCES users(id)
                    ON DELETE CASCADE )`
	DeleteData       = "DELETE FROM storage WHERE uid = $1 AND id = $2"
	GetAllDataByType = "SELECT * FROM storage WHERE uid = $1 AND type = $2"
	GetDataByID      = "SELECT * FROM storage WHERE uid = $1 AND id = $2"
	StoreData        = `
		INSERT INTO storage(uid, data, type) VALUES($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id
	`
)

func NewDBRepo(url string) (*DBRepo, error) {
	if url == "" {
		return &DBRepo{}, ErrDBMissingURL
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &DBRepo{}, err
	}

	_, err = db.ExecContext(context.Background(), CreateStorageTable)
	return &DBRepo{db: db}, err
}

func (r *DBRepo) DeleteData(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrNotFound
	}

	res, err := r.db.ExecContext(ctx, DeleteData, uid, id)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *DBRepo) GetAllDataByType(ctx context.Context, uid string,
	t StorageType,
) ([]SecureData, error) {
	if uid == "" {
		return nil, ErrMissingArgs
	}

	rows, err := r.db.QueryContext(ctx, GetAllDataByType, uid, t)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer r.closeRows(rows)

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
	if uid == "" || id == "" {
		return SecureData{}, ErrNotFound
	}

	var data SecureData
	err := r.db.QueryRowContext(ctx, GetDataByID, uid, id).Scan(&data.ID, &data.UID, &data.Data, &data.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return SecureData{}, ErrNotFound
	}
	return data, err
}

func (r *DBRepo) StoreData(ctx context.Context, data SecureData) (string, error) {
	if data.Data == nil || data.UID == "" {
		return "", ErrEmpty
	}

	var id string
	err := r.db.QueryRowContext(ctx, StoreData, data.UID, data.Data, data.Type).Scan(&id)
	return id, err
}

func (r *DBRepo) closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Error(err)
	}
}
