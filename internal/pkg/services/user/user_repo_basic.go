package user

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BasicRepo struct {
	users *sync.Map
}

func NewBasicRepo() *BasicRepo {
	return &BasicRepo{users: &sync.Map{}}
}

func (r *BasicRepo) AddUser(_ context.Context, user User) (User, error) {
	if user.Name == "" || user.Password == "" {
		return User{}, ErrCredMissing
	}

	if _, ok := r.users.Load(user.ID); user.ID != "" && ok {
		return User{}, ErrExists
	}

	id := uuid.NewString()
	user.ID = id
	r.users.Store(id, user)
	return user, nil
}

func (r *BasicRepo) DeleteUser(_ context.Context, uid string) error {
	if _, ok := r.users.Load(uid); !ok || uid == "" {
		return ErrNotFound
	}
	r.users.Delete(uid)
	return nil
}

func (r *BasicRepo) GetUserByID(_ context.Context, uid string) (User, error) {
	if uid == "" {
		return User{}, ErrNotFound
	}
	if u, ok := r.users.Load(uid); ok {
		return u.(User), nil
	}
	return User{}, ErrNotFound
}

func (r *BasicRepo) GetUserByName(_ context.Context, name string) (User, error) {
	if name == "" {
		return User{}, ErrNotFound
	}

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
