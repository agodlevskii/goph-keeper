package session

import (
	"context"
	"sync"
)

type BasicRepo struct {
	tokens sync.Map
}

func NewBasicRepo() *BasicRepo {
	return &BasicRepo{tokens: sync.Map{}}
}

func (r *BasicRepo) DeleteSession(_ context.Context, cid string) error {
	r.tokens.Delete(cid)
	return nil
}

func (r *BasicRepo) GetSession(_ context.Context, cid string) (string, error) {
	if t, ok := r.tokens.Load(cid); ok {
		return t.(string), nil
	}
	return "", ErrNotFound
}

func (r *BasicRepo) StoreSession(_ context.Context, cid, token string) error {
	r.tokens.Store(cid, token)
	return nil
}
