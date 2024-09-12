package state_store

import (
	"fmt"

	state_types "github.com/errorboi/boo/types/user_state"
)

type StateStore interface {
	Get(userID int64) (*state_types.State, error)
	Init(userID int64) (*state_types.State, error)
	SetState(userID int64, state *state_types.State) error
	Del(userID int64) error
}

type inmemStateStore struct {
	prefix string
	store  map[string]*state_types.State
}

func NewInmemStateStore(storeName string) StateStore {
	return &inmemStateStore{
		store:  make(map[string]*state_types.State),
		prefix: storeName,
	}
}

func (s *inmemStateStore) Get(userID int64) (*state_types.State, error) {
	step, ok := s.store[s.key(userID)]
	if !ok {
		return nil, state_types.ErrStepNotFound
	}
	return step, nil
}

func (s *inmemStateStore) Init(userID int64) (*state_types.State, error) {
	s.store[s.key(userID)] = state_types.NewState(state_types.Init)
	return s.store[s.key(userID)], nil
}

func (s *inmemStateStore) SetState(userID int64, state *state_types.State) error {
	s.store[s.key(userID)] = state
	return nil
}

func (s *inmemStateStore) Del(userID int64) error {
	delete(s.store, s.key(userID))
	return nil

}

func (s *inmemStateStore) key(userID int64) string {
	return fmt.Sprintf("%s/%d", s.prefix, userID)
}
