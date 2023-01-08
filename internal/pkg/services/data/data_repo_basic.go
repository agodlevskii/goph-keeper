package data

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BasicRepo struct {
	data *sync.Map
}

type Storage struct {
	user *sync.Map
}

func NewBasicRepo() *BasicRepo {
	return &BasicRepo{data: &sync.Map{}}
}

func (r *BasicRepo) DeleteData(_ context.Context, uid, id string) error {
	if us, ok := r.data.Load(uid); ok {
		if _, ok = us.(Storage).user.Load(id); ok {
			us.(Storage).user.Delete(id)
			return nil
		}
	}
	return ErrNotFound
}

func (r *BasicRepo) GetAllDataByType(_ context.Context, uid string,
	t StorageType,
) ([]SecureData, error) {
	if uid == "" {
		return nil, ErrMissingArgs
	}

	var data []SecureData
	if us, ok := r.data.Load(uid); ok {
		us.(Storage).user.Range(func(_, v any) bool {
			d := v.(SecureData)
			if d.Type == t {
				data = append(data, d)
			}
			return true
		})
	}
	return data, nil
}

func (r *BasicRepo) GetDataByID(_ context.Context, uid, id string) (SecureData, error) {
	var (
		us any
		d  any
		ok bool
	)

	if us, ok = r.data.Load(uid); ok {
		if d, ok = us.(Storage).user.Load(id); ok {
			return d.(SecureData), nil
		}
	}
	return SecureData{}, ErrNotFound
}

func (r *BasicRepo) StoreData(_ context.Context, data SecureData) (string, error) {
	if data.Data == nil || data.UID == "" {
		return "", ErrEmpty
	}

	id := uuid.NewString()
	data.ID = id

	if us, ok := r.data.Load(data.UID); !ok {
		sd := &sync.Map{}
		sd.Store(id, data)
		r.data.Store(data.UID, Storage{user: sd})
	} else {
		us.(Storage).user.Store(id, data)
	}

	return id, nil
}
