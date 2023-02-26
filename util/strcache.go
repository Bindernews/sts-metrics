package util

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
type StrStoreFn func(context.Context, []string) error

// Function that takes a list of strings and resolves them into IDs.
type StrLoadFn func(context.Context, []string) ([]int32, error)

type StrCache interface {
	Load(ctx context.Context, strings ...[]string) error
	Get(key string) int32
	MaybeGet(key *string) sql.NullInt32
	GetAll(keys []string) []int32
}

// Cache of strings to StrCache.id
type syncStrCache struct {
	strs    map[string]int32
	lock    sync.RWMutex
	StoreFn StrStoreFn
	LoadFn  StrLoadFn
}

func NewStrCache(loadFn StrLoadFn, storeFn StrStoreFn) StrCache {
	c := syncStrCache{
		strs:    make(map[string]int32),
		StoreFn: storeFn,
		LoadFn:  loadFn,
	}
	return &c
}

// Loads all strings used by the passed users
func (s *syncStrCache) Load(ctx context.Context, strings ...[]string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	needed := make(map[string]int32)
	// Gather all needed strings
	for _, u := range strings {
		for _, v := range u {
			if _, ok := s.strs[v]; !ok {
				needed[v] = 0
			}
		}
	}
	keys := lo.Keys(needed)
	var values []int32
	// Query from DB
	if s.StoreFn != nil {
		if err = s.StoreFn(ctx, keys); err != nil {
			log.Println(err)
			return ErrCouldNotStoreStrings
		}
	}
	if values, err = s.LoadFn(ctx, keys); err != nil {
		return
	}
	if len(values) != len(keys) {
		return ErrCouldNotStoreStrings
	}
	for i, v := range values {
		s.strs[keys[i]] = v
	}
	return nil
}

func (s *syncStrCache) Get(k string) int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.strs[k]
}

func (s *syncStrCache) MaybeGet(k *string) sql.NullInt32 {
	if k == nil {
		return sql.NullInt32{Valid: false}
	} else {
		s.lock.RLock()
		defer s.lock.RUnlock()
		return sql.NullInt32{Valid: true, Int32: s.Get(*k)}
	}
}

func (s *syncStrCache) GetAll(keys []string) []int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return lo.Map(keys, func(k string, _ int) int32 { return s.strs[k] })
}
