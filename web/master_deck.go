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
	Card orm.CardSpec
	// Instances in master deck
	Count int
}

func (e DeckEntry) ToOrm(oc *OrmContext, _ int) orm.AddMasterDeckParams {
	return orm.AddMasterDeckParams{
		RunID:  oc.Runid,
		CardID: oc.Cc.Get(e.Card),
		Count:  int16(e.Count),
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
			row = DeckEntry{Card: CardNameSplit(card), Count: 1}
		}
		(*d)[card] = row
	}
}

func (d MasterDeck) ToOrm(oc *OrmContext) []orm.AddMasterDeckParams {
	return MapToOrm[orm.AddMasterDeckParams](lo.Values(d), oc)
}

func (d MasterDeck) GetCards() []orm.CardSpec {
	cards := make([]orm.CardSpec, 0)
	for _, e := range d {
		cards = append(cards, e.Card)
	}
	return cards
}

// Returns list of unique base names in the deck
func (d MasterDeck) GetStrings() []string {
	names := make(map[string]bool)
	for _, e := range d {
		names[e.Card.Card] = true
	}
	return lo.Keys(names)
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
