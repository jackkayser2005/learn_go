package store

import (
	"sync"
)

type UserLink struct {
	Text string `json:"text"`
}

type MemStore struct {
	mux sync.RWMutex
	m   map[string]UserLink
}

func NewMemStore() *MemStore {
	return &MemStore{m: make(map[string]UserLink)}
}

func (s *MemStore) Save(code string, link UserLink) {
	s.mux.Lock()
	s.m[code] = link
	s.mux.Unlock()
}

func (s *MemStore) Find(code string) (UserLink, bool) {
	s.mux.RLock()
	link, ok := s.m[code]
	s.mux.RUnlock()
	return link, ok
}
