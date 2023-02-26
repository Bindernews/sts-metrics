package stms

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strconv"

	"github.com/bindernews/sts-msr/orm"
	"github.com/samber/lo"
)

// Regex that splits a card into base name and upgrade count.
// If the card name does not match, it has not been upgraded.
var CardUpgradeRegex = regexp.MustCompile(`(.+)\+([0-9]+)$`)

// Error when StrCache can't resolve all strings
var ErrCouldNotLoadStrings = errors.New("could not load strings")

var ErrCouldNotStoreStrings = errors.New("cannot add new strings")

// Function that adds a list of strings to the cache.
type StrStoreFn func(context.Context, []string) error

// Function that takes a list of strings and resolves them into IDs.
type StrLoadFn func(context.Context, []string) ([]int32, error)

// Cache of strings to StrCache.id
type StrCache struct {
	strs    map[string]int32
	StoreFn StrStoreFn
	LoadFn  StrLoadFn
}

func NewStrCache(loadFn StrLoadFn, storeFn StrStoreFn) *StrCache {
	c := StrCache{
		strs:    make(map[string]int32),
		StoreFn: storeFn,
		LoadFn:  loadFn,
	}
	return &c
}

// Loads all strings used by the passed users
func (s *StrCache) Load(ctx context.Context, strings ...[]string) (err error) {
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

func (s *StrCache) Get(k string) int32 {
	return s.strs[k]
}

func (s *StrCache) MaybeGet(k *string) sql.NullInt32 {
	if k == nil {
		return sql.NullInt32{Valid: false}
	} else {
		return sql.NullInt32{Valid: true, Int32: s.Get(*k)}
	}
}

func (s *StrCache) GetAll(keys []string) []int32 {
	return lo.Map(keys, func(k string, _ int) int32 { return s.strs[k] })
}

type StrCacheUser interface {
	// Returns the list of all cacheable strings used by this object
	GetStrings() []string
}

type ConvToOrm[T any] interface {
	ToOrm(sc *StrCache, runid int32) T
}

func MapToOrm[T any, E ConvToOrm[T]](inp []E, sc *StrCache, runid int32) []T {
	out := make([]T, len(inp))
	for i, e := range inp {
		out[i] = e.ToOrm(sc, runid)
	}
	return out
}

// Parsed master deck entry
type DeckEntry struct {
	// Base card name (without upgrades)
	Name string
	// Upgrade count
	Upgrades int
	// Instances in master deck
	Count int
}

func (e DeckEntry) ToOrm(sc *StrCache, runid int32) orm.AddMasterDeckParams {
	return orm.AddMasterDeckParams{
		RunID:    runid,
		CardID:   sc.Get(e.Name),
		Count:    int16(e.Count),
		Upgrades: int16(e.Upgrades),
	}
}

// Master deck represented as map of card names to entries
type MasterDeck map[string]DeckEntry

func NewMasterDeck() *MasterDeck {
	d := new(MasterDeck)
	*d = make(map[string]DeckEntry)
	return d
}

// Add an array of card names (which may be upgraded) to this deck.
func (d *MasterDeck) AddCards(cards []string) {
	var row DeckEntry
	var ok bool
	for _, card := range cards {
		// If we've already seen this card, increment the count
		if row, ok = (*d)[card]; ok {
			row.Count += 1
		} else {
			// Otherwise add a new entry for it
			name, upg := CardNameSplit(card)
			row = DeckEntry{Name: name, Upgrades: upg, Count: 1}
		}
		(*d)[card] = row
	}
}

// Returns list of unique base names in the deck
func (d MasterDeck) GetAllNames() []string {
	names := make(map[string]bool)
	for _, e := range d {
		names[e.Name] = true
	}
	return lo.Keys(names)
}

func (d MasterDeck) ToOrm(sc *StrCache, runid int32) []orm.AddMasterDeckParams {
	return MapToOrm[orm.AddMasterDeckParams](lo.Values(d), sc, runid)
}

func (d MasterDeck) GetStrings() []string {
	return d.GetAllNames()
}

// Takes a card name that may include an upgrade count
// and returns the base name and number of ugrades (may be 0).
func CardNameSplit(card string) (string, int) {
	if mt := CardUpgradeRegex.FindStringSubmatch(card); mt != nil {
		v, err := strconv.Atoi(mt[2])
		if err != nil {
			v = 0
		}
		return mt[1], v
	} else {
		return card, 0
	}
}
