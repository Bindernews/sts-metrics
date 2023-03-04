package web

import (
	"regexp"
	"strconv"

	"github.com/bindernews/sts-msr/orm"
	"github.com/samber/lo"
)

// Regex that splits a card into base name and upgrade count.
// If the card name does not match, it has not been upgraded.
var CardUpgradeRegex = regexp.MustCompile(`(.+)\+([0-9]+)$`)

// Parsed master deck entry
type DeckEntry struct {
	// Base card name (without upgrades)
	Name string
	// Upgrade count
	Upgrades int
	// Instances in master deck
	Count int
}

func (e DeckEntry) ToOrm(sc StrCache, runid int32) orm.AddMasterDeckParams {
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

func (d MasterDeck) ToOrm(sc StrCache, runid int32) []orm.AddMasterDeckParams {
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
