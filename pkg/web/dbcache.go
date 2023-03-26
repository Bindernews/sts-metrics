package web

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/samber/lo"
)

// Error when StrCache can't resolve all strings
var ErrCouldNotLoadStrings = errors.New("could not load strings")

var ErrCouldNotStoreStrings = errors.New("cannot add new strings")

// Function that adds a list of strings to the cache.
type DbStoreFn[T comparable] func(context.Context, []T) error

// Function that takes a list of strings and resolves them into IDs.
type DbLoadFn[T comparable] func(context.Context, []T) ([]int32, error)

// Cache of T values to their database ID
type DbCache[T comparable] interface {
	Load(ctx context.Context, values ...[]T) error
	Get(key T) int32
	MaybeGet(key *T) sql.NullInt32
	GetAll(keys []T) []int32
	RevGet(id int32) T
	RevGetAll(id []int32) []T
}

type syncDbCache[T comparable] struct {
	fwd     map[T]int32
	rev     map[int32]T
	lock    sync.RWMutex
	storeFn DbStoreFn[T]
	loadFn  DbLoadFn[T]
}

func NewDbCache[T comparable](loadFn DbLoadFn[T], storeFn DbStoreFn[T]) DbCache[T] {
	c := syncDbCache[T]{
		fwd:     make(map[T]int32),
		rev:     make(map[int32]T),
		storeFn: storeFn,
		loadFn:  loadFn,
	}
	return &c
}

// Loads all strings used by the passed users
func (s *syncDbCache[T]) Load(ctx context.Context, values ...[]T) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	needed := make(map[T]int32)
	// Gather all needed strings
	for _, u := range values {
		for _, v := range u {
			if _, ok := s.fwd[v]; !ok {
				needed[v] = 0
			}
		}
	}
	keys := lo.Keys(needed)
	var dbIds []int32
	// Query from DB
	if s.storeFn != nil {
		if err = s.storeFn(ctx, keys); err != nil {
			log.Println(err)
			return ErrCouldNotStoreStrings
		}
	}
	if dbIds, err = s.loadFn(ctx, keys); err != nil {
		return
	}
	if len(dbIds) != len(keys) {
		return ErrCouldNotStoreStrings
	}
	for i, v := range dbIds {
		s.fwd[keys[i]] = v
		s.rev[v] = keys[i]
	}

	return nil
}

func (s *syncDbCache[T]) Get(k T) int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.fwd[k]
}

func (s *syncDbCache[T]) MaybeGet(k *T) sql.NullInt32 {
	if k == nil {
		return sql.NullInt32{Valid: false}
	} else {
		s.lock.RLock()
		defer s.lock.RUnlock()
		return sql.NullInt32{Valid: true, Int32: s.Get(*k)}
	}
}

func (s *syncDbCache[T]) GetAll(keys []T) []int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return lo.Map(keys, func(k T, _ int) int32 { return s.fwd[k] })
}

func (s *syncDbCache[T]) RevGet(id int32) T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.rev[id]
}

func (s *syncDbCache[T]) RevGetAll(keys []int32) []T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	out := make([]T, len(keys))
	for i, k := range keys {
		out[i] = s.rev[k]
	}
	return out
}
