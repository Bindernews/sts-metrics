package web

import (
	"regexp"
	"strconv"

	"github.com/bindernews/sts-msr/orm"
)

// Regex that splits a card into base name and upgrade count.
// If the card name does not match, it has not been upgraded.
var CardUpgradeRegex = regexp.MustCompile(`(.+)\+([0-9]+)$`)

// Caches and other information needed to convert JSON data to DB structures.
type OrmContext struct {
	// String cache
	Sc DbCache[string]
	// Card cache
	Cc DbCache[orm.CardSpec]
	// Set of strings that this run will need to load
	StringSet map[string]bool
	// Set of cards this run will need to load
	CardSet map[orm.CardSpec]bool
	// ID of run in database
	Runid int32
}

type ConvToOrm interface {
	// Convert the item to its orm type T. ix is the array index of the item.
	ToOrm(oc *OrmContext, ix int) any
}

type IOrmPreload interface {
	// Give the item the opportunity to add things to the list of items
	// to be cached. Caches are updated in bulk to reduce database IO.
	Preload(oc *OrmContext, ix int)
}

func PreloadArray[T IOrmPreload](oc *OrmContext, inp []T) {
	for i, v := range inp {
		v.Preload(oc, i)
	}
}

// Calls ToOrm on inp, returning the result
func MapToOrm[T any](oc *OrmContext, inp []ConvToOrm) []T {
	out := make([]T, len(inp))
	for i, e := range inp {
		out[i] = e.ToOrm(oc, i).(T)
	}
	return out
}

// Makes a copy of the OrmContext with per-run data reset
func (oc OrmContext) Copy() *OrmContext {
	return &OrmContext{
		Sc:        oc.Sc,
		Cc:        oc.Cc,
		StringSet: make(map[string]bool),
		CardSet:   make(map[orm.CardSpec]bool),
		Runid:     0,
	}
}

// Takes a card name that may include an upgrade count
// and returns the base name and number of ugrades (may be 0).
func CardNameSplit(card string) orm.CardSpec {
	if mt := CardUpgradeRegex.FindStringSubmatch(card); mt != nil {
		v, err := strconv.Atoi(mt[2])
		if err != nil {
			v = 0
		}
		return orm.CardSpec{Card: mt[1], Upgrades: v}
	} else {
		return orm.CardSpec{Card: card, Upgrades: 0}
	}
}

func StringsToCards(cards []string) []orm.CardSpec {
	out := make([]orm.CardSpec, len(cards))
	for i, v := range cards {
		out[i] = CardNameSplit(v)
	}
	return out
}
