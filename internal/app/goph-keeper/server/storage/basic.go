package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type BasicRepo struct {
	data   sync.Map
	tokens sync.Map
	users  sync.Map
}

type DataStorage struct {
	user *sync.Map
}

func NewBasicStorage() *BasicRepo {
	return &BasicRepo{
		data:   sync.Map{},
		tokens: sync.Map{},
		users:  sync.Map{},
	}
}

func (r *BasicRepo) DeleteSession(_ context.Context, cid string) error {
	r.tokens.Delete(cid)
	return nil
}

func (r *BasicRepo) GetSession(_ context.Context, cid string) (string, error) {
	if t, ok := r.tokens.Load(cid); ok {
		return t.(string), nil
	}
	return "", errors.New("token not found")
}

func (r *BasicRepo) StoreSession(_ context.Context, cid, token string) error {
	r.tokens.Store(cid, token)
	return nil
}

func (r *BasicRepo) AddUser(_ context.Context, u User) (User, error) {
	id := uuid.NewString()
	u.ID = id
	r.users.Store(id, u)
	return u, nil
}

func (r *BasicRepo) GetUserByID(_ context.Context, uid string) (User, error) {
	if u, ok := r.users.Load(uid); ok {
		return u.(User), nil
	}
	return User{}, errors.New("user not found")
}

func (r *BasicRepo) GetUserByName(_ context.Context, name string) (User, error) {
	var user User

	r.users.Range(func(_, v any) bool {
		u := v.(User)
		if u.Name == name {
			user = u
			return false
		}
		return true
	})

	if user.ID == "" {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *BasicRepo) GetAllDataByType(_ context.Context, uid string, t Type) ([]SecureData, error) {
	var data []SecureData
	if us, ok := r.data.Load(uid); ok {
		us.(DataStorage).user.Range(func(_, v any) bool {
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
		if d, ok = us.(DataStorage).user.Load(id); ok {
			data := d.(SecureData)
			return data, nil
		}
	}
	return SecureData{}, errors.New("data not found")
}

func (r *BasicRepo) StoreData(_ context.Context, data SecureData) (string, error) {
	id := uuid.NewString()
	data.ID = id

	if us, ok := r.data.Load(data.UID); !ok {
		sd := &sync.Map{}
		sd.Store(id, data)
		r.data.Store(data.UID, DataStorage{user: sd})
	} else {
		us.(DataStorage).user.Store(id, data)
	}

	return id, nil
}
