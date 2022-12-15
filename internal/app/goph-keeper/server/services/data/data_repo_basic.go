package data

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BasicRepo struct {
	data sync.Map
}

type Storage struct {
	user *sync.Map
}

func NewBasicRepo() *BasicRepo {
	return &BasicRepo{data: sync.Map{}}
}

func (r *BasicRepo) DeleteData(_ context.Context, uid, id string) error {
	if d, ok := r.data.Load(id); ok {
		data := d.(SecureData)
		if data.UID == uid {
			r.data.Delete(id)
		}
	}
	return nil
}

func (r *BasicRepo) GetAllDataByType(_ context.Context, uid string,
	t StorageType,
) ([]SecureData, error) {
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
			data := d.(SecureData)
			return data, nil
		}
	}
	return SecureData{}, ErrNotFound
}

func (r *BasicRepo) StoreData(_ context.Context, data SecureData) (string, error) {
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
