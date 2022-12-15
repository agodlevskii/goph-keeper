package user

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BasicRepo struct {
	users sync.Map
}

func NewBasicRepo() *BasicRepo {
	return &BasicRepo{users: sync.Map{}}
}

func (r *BasicRepo) AddUser(_ context.Context, user User) (User, error) {
	id := uuid.NewString()
	user.ID = id
	r.users.Store(id, user)
	return user, nil
}

func (r *BasicRepo) DeleteUser(_ context.Context, uid string) error {
	r.users.Delete(uid)
	return nil
}

func (r *BasicRepo) GetUserByID(_ context.Context, uid string) (User, error) {
	if u, ok := r.users.Load(uid); ok {
		return u.(User), nil
	}
	return User{}, ErrNotFound
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
		return User{}, ErrNotFound
	}
	return user, nil
}
