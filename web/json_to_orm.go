package web

// Helpers for mapping the generated JSON structs to the generated SQL structs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bindernews/sts-msr/orm"
	"github.com/jackc/pgtype"
	"github.com/samber/lo"
)

type OrmContext struct {
	// String cache
	Sc DbCache[string]
	// Card cache
	Cc DbCache[orm.CardSpec]
	// ID of run in database
	Runid int32
}

func (oc OrmContext) Copy() *OrmContext {
	return &OrmContext{
		Sc:    oc.Sc,
		Cc:    oc.Cc,
		Runid: 0,
	}
}

type ConvToOrm[T any] interface {
	ToOrm(oc *OrmContext, ix int) T
}

// Calls ToOrm on inp, returning the result
func MapToOrm[T any, E ConvToOrm[T]](inp []E, oc *OrmContext) []T {
	out := make([]T, len(inp))
	for i, e := range inp {
		out[i] = e.ToOrm(oc, i)
	}
	return out
}

// Make sure data is properly round-tripped.
// It's not the fastest method, but oh well.
func (j *RunSchemaJson) MarshalJSON() ([]byte, error) {
	// NOTE: could use mapstucture library, but again, perf is not that important here
	type Plain RunSchemaJson
	plain := Plain(*j)
	plainData, err := json.Marshal(plain)
	if err != nil {
		return nil, err
	}
	raw := make(map[string]any)
	if err := json.Unmarshal(plainData, &raw); err != nil {
		return nil, err
	}
	for k, v := range j.Extra {
		raw[k] = v
	}
	return json.Marshal(raw)
}

func (c BossRelicChoice) ToOrm(oc *OrmContext, ix int) orm.AddBossRelicsParams {
	return orm.AddBossRelicsParams{
		RunID:     oc.Runid,
		NotPicked: oc.Sc.GetAll(c.NotPicked),
		Picked:    oc.Sc.Get(c.Picked),
		Ord:       int16(ix),
	}
}

func (c *BossRelicChoice) GetStrings() []string {
	return append([]string{c.Picked}, c.NotPicked...)
}

// Convert CampfireChoice to an orm object
func (c CampfireChoice) ToOrm(oc *OrmContext, _ int) orm.AddCampfireParams {
	return orm.AddCampfireParams{
		RunID: oc.Runid,
		Data:  oc.Sc.MaybeGet(c.Data),
		Key:   oc.Sc.Get(c.Key),
		Floor: int32(c.Floor),
	}
}

func (c *CampfireChoice) GetStrings() []string {
	data := ""
	if c.Data != nil {
		data = *c.Data
	}
	return []string{c.Key, data}
}

// Convert CardChoice to intermediary form
func (c CardChoice) ToOrm(oc *OrmContext, _ int) CardChoiceParsed {
	cp := CardChoiceParsed{
		Floor:     int(c.Floor),
		NotPicked: make([]orm.CardSpec, len(c.NotPicked)),
		Picked:    CardNameSplit(c.Picked),
	}
	for i, v := range c.NotPicked {
		cp.NotPicked[i] = CardNameSplit(v)
	}
	return cp
}

// Pre-parsed version of CardChoice
type CardChoiceParsed struct {
	Floor     int
	NotPicked []orm.CardSpec
	Picked    orm.CardSpec
}

func (cp CardChoiceParsed) GetCards() []orm.CardSpec {
	return append([]orm.CardSpec{cp.Picked}, cp.NotPicked...)
}

func (cp CardChoiceParsed) ToOrm(oc *OrmContext, _ int) orm.AddCardChoiceParams {
	return orm.AddCardChoiceParams{
		RunID:     oc.Runid,
		Floor:     int32(cp.Floor),
		NotPicked: oc.Cc.GetAll(cp.NotPicked),
		Picked:    oc.Cc.Get(cp.Picked),
	}
}

func (c DamageTaken) ToOrm(oc *OrmContext, _ int) orm.AddDamageTakenParams {
	return orm.AddDamageTakenParams{
		RunID:   oc.Runid,
		Floor:   int32(c.Floor),
		Turns:   int32(c.Turns),
		Damage:  float32(c.Damage),
		Enemies: oc.Sc.Get(c.Enemies),
	}
}

func (c *DamageTaken) GetStrings() []string {
	return []string{c.Enemies}
}

func (c EventChoice) ToOrm(oc *OrmContext, _ int) orm.AddEventChoicesParams {
	return orm.AddEventChoicesParams{
		RunID:             oc.Runid,
		DamageDelta:       int32(c.DamageHealed - c.DamageTaken),
		EventNameID:       oc.Sc.Get(c.EventName),
		Floor:             int32(c.Floor),
		GoldDelta:         int32(c.GoldGain - c.GoldLoss),
		MaxHpDelta:        int32(c.MaxHpGain - c.MaxHpLoss),
		PlayerChoiceID:    oc.Sc.Get(c.PlayerChoice),
		RelicsObtainedIds: oc.Sc.GetAll(c.RelicsObtained),
	}
}

func (c *EventChoice) GetStrings() []string {
	return append([]string{c.EventName, c.PlayerChoice}, c.RelicsObtained...)
}

func (c PotionObtained) ToOrm(oc *OrmContext, _ int) orm.AddPotionObtainParams {
	return orm.AddPotionObtainParams{
		RunID: oc.Runid,
		Floor: int16(c.Floor),
		Key:   oc.Sc.Get(c.Key),
	}
}

func (c *PotionObtained) GetStrings() []string {
	return []string{c.Key}
}

func (s RelicObtain) ToOrm(oc *OrmContext, _ int) orm.AddRelicObtainParams {
	return orm.AddRelicObtainParams{
		RunID: oc.Runid,
		Floor: int16(s.Floor),
		Key:   oc.Sc.Get(s.Key),
	}
}

func (c *RelicObtain) GetStrings() []string {
	return []string{c.Key}
}

func (s *RunSchemaJson) ToAddRunRaw(oc *OrmContext) orm.AddRunRawParams {
	tstamp := time.UnixMilli(int64(s.Timestamp))
	pathNorm := lo.Map(s.PathPerFloor, func(v FloorPath, _ int) string {
		return DeNull(v)
	})

	return orm.AddRunRawParams{
		AscensionLevel:   int32(s.AscensionLevel),
		BuildVersion:     oc.Sc.Get(s.BuildVersion),
		CampfireRested:   int32(s.CampfireRested),
		CampfireUpgraded: int32(s.CampfireUpgraded),
		CharacterID:      oc.Sc.Get(s.CharacterChosen),
		ChooseSeed:       s.ChoseSeed,
		CircletCount:     int32(s.CircletCount),
		FloorReached:     int32(s.FloorReached),
		Gold:             int32(s.Gold),
		KilledBy:         oc.Sc.Get(s.KilledBy),
		LocalTime:        s.LocalTime,
		NeowBonusID:      oc.Sc.Get(s.NeowBonus),
		NeowCostID:       oc.Sc.Get(s.NeowCost),
		PathPerFloor:     pathToStringFwd(pathNorm),
		PathTaken:        pathToStringFwd(s.PathTaken),
		PlayID:           s.PlayId.String(),
		PlayerExperience: int32(s.PlayerExperience),
		Playtime:         int32(s.Playtime),
		PurchasedPurges:  int32(s.PurchasedPurges),
		Score:            int32(s.Score),
		SeedPlayed:       s.SeedPlayed,
		//
		SeedSourceTimestamp: sqlInt32(s.SeedSourceTimestamp),
		SpecialSeed:         int32(s.SpecialSeed),
		Timestamp:           sql.NullTime{Time: tstamp, Valid: true},
		Victory:             s.Victory,
		WinRate:             s.WinRate,
	}
}

func (s *RunSchemaJson) toPerFloorOrm(oc *OrmContext, runid int32) []orm.AddPerFloorParams {
	end := len(s.CurrentHpPerFloor)
	out := make([]orm.AddPerFloorParams, end)
	for i := 0; i < end; i++ {
		out[i] = orm.AddPerFloorParams{
			RunID:     runid,
			Floor:     int16(i),
			Gold:      int32(s.GoldPerFloor[i]),
			CurrentHp: int32(s.CurrentHpPerFloor[i]),
			MaxHp:     int32(s.MaxHpPerFloor[i]),
		}
	}
	return out
}

func (s RunSchemaJson) getOtherStrings() []string {
	out := []string{}
	out = append(out, s.DailyMods...)
	out = append(out, s.ItemsPurchased...)
	out = append(out, s.ItemsPurged...)
	out = append(out, s.Relics...)
	for _, u := range s.BossRelics {
		out = append(out, u.GetStrings()...)
	}
	for _, u := range s.CampfireChoices {
		out = append(out, u.GetStrings()...)
	}
	for _, u := range s.DamageTaken {
		out = append(out, u.GetStrings()...)
	}
	for _, u := range s.EventChoices {
		out = append(out, u.GetStrings()...)
	}
	for _, u := range s.PotionsObtained {
		out = append(out, u.GetStrings()...)
	}
	for _, u := range s.RelicsObtained {
		out = append(out, u.GetStrings()...)
	}
	return out
}

func (s RunSchemaJson) getMinimalStrings() []string {
	return []string{s.BuildVersion, s.CharacterChosen, s.KilledBy, s.NeowBonus, s.NeowCost}
}

// Add this Run to the database. Returns the rowid of the run.
func (r *RunSchemaJson) AddToDb(ctx context.Context, oc *OrmContext, db *orm.Queries) (runId int32, err error) {
	// Pre-load a few strings so we can add the run with valid references
	// Note that this MAY cause us to add strings we don't use if the play_id is a duplicate.
	if err = oc.Sc.Load(ctx, r.getMinimalStrings()); err != nil {
		return
	}
	// Insert the new run. May fail if the run already exists.
	runId, err = db.AddRunRaw(ctx, r.ToAddRunRaw(oc))
	if err != nil {
		return
	}
	oc.Runid = runId

	// Cache other strings
	if err = oc.Sc.Load(ctx, r.getOtherStrings()); err != nil {
		return
	}

	// Parse deck and other card information so we can cache CardSpecs
	deck := NewMasterDeck()
	deck.AddCards(r.MasterDeck)
	// Parse items purchased + purged
	specsPurchased := StringsToCards(r.ItemsPurchased)
	specsPurged := StringsToCards(r.ItemsPurged)
	// Parse card choices
	parsedCards := MapToOrm[CardChoiceParsed](r.CardChoices, oc)
	ormCardSpecs := lo.FlatMap(parsedCards, func(cp CardChoiceParsed, _ int) []orm.CardSpec {
		return cp.GetCards()
	})
	// Cache CardSpecs
	if err = oc.Cc.Load(ctx, ormCardSpecs, deck.GetCards(), specsPurchased, specsPurged); err != nil {
		return
	}

	// Add items purchased
	itemsPurchased := make([]orm.AddItemsPurchasedParams, len(specsPurchased))
	for i, v := range r.ItemPurchaseFloors {
		itemsPurchased[i] = orm.AddItemsPurchasedParams{
			RunID:  oc.Runid,
			CardID: oc.Cc.Get(specsPurchased[i]),
			Floor:  int16(v),
		}
	}
	if _, err = db.AddItemsPurchased(ctx, itemsPurchased); err != nil {
		return
	}
	// Add items purged
	itemsPurged := make([]orm.AddItemsPurgedParams, len(specsPurged))
	for i, v := range r.ItemsPurgedFloors {
		itemsPurged[i] = orm.AddItemsPurgedParams{
			RunID:  oc.Runid,
			CardID: oc.Cc.Get(specsPurged[i]),
			Floor:  int16(v),
		}
	}
	if _, err = db.AddItemsPurged(ctx, itemsPurged); err != nil {
		return
	}
	// Add card choices
	cardChoices := MapToOrm[orm.AddCardChoiceParams](parsedCards, oc)
	if _, err = db.AddCardChoice(ctx, cardChoices); err != nil {
		return
	}
	// Add arrays data
	ormArrays := []orm.AddRunArraysParams{{
		RunID:               oc.Runid,
		DailyMods:           oc.Sc.GetAll(r.DailyMods),
		PotionsFloorSpawned: mapInt32(r.PotionsFloorSpawned),
		PotionsFloorUsage:   mapInt32(r.PotionsFloorUsage),
		RelicIds:            oc.Sc.GetAll(r.Relics),
	}}
	if _, err = db.AddRunArrays(ctx, ormArrays); err != nil {
		return
	}
	if _, err = db.AddMasterDeck(ctx, deck.ToOrm(oc)); err != nil {
		return
	}

	// Flags
	flags := map[string]bool{
		"ascension": r.IsAscensionMode,
		"beta":      r.IsBeta,
		"daily":     r.IsDaily,
		"endless":   r.IsEndless,
		"prod":      r.IsProd,
		"trial":     r.IsTrial,
	}
	for k, ok := range flags {
		if ok {
			if err = db.AddFlag(ctx, orm.AddFlagParams{
				RunID: runId,
				Flag:  orm.FlagKind(k),
			}); err != nil {
				return
			}
		}
	}

	// Store unparsed data
	if len(r.Extra) > 0 {
		var extraBytes []byte
		if extraBytes, err = json.Marshal(r.Extra); err != nil {
			return
		}
		if err = db.AddRunsExtra(ctx, orm.AddRunsExtraParams{
			RunID: runId,
			Extra: pgtype.JSONB{Bytes: extraBytes, Status: pgtype.Present},
		}); err != nil {
			return
		}
	}
	// All other DB rows
	perFloor := r.toPerFloorOrm(oc, runId)
	if _, err = db.AddPerFloor(ctx, perFloor); err != nil {
		return
	}
	ormBossRelics := MapToOrm[orm.AddBossRelicsParams](r.BossRelics, oc)
	if _, err = db.AddBossRelics(ctx, ormBossRelics); err != nil {
		return
	}
	ormCampfires := MapToOrm[orm.AddCampfireParams](r.CampfireChoices, oc)
	if _, err = db.AddCampfire(ctx, ormCampfires); err != nil {
		return
	}
	ormDamageTaken := MapToOrm[orm.AddDamageTakenParams](r.DamageTaken, oc)
	if _, err = db.AddDamageTaken(ctx, ormDamageTaken); err != nil {
		return
	}
	ormEvents := MapToOrm[orm.AddEventChoicesParams](r.EventChoices, oc)
	if _, err = db.AddEventChoices(ctx, ormEvents); err != nil {
		return
	}
	ormPotions := MapToOrm[orm.AddPotionObtainParams](r.PotionsObtained, oc)
	if _, err = db.AddPotionObtain(ctx, ormPotions); err != nil {
		return
	}
	ormRelics := MapToOrm[orm.AddRelicObtainParams](r.RelicsObtained, oc)
	if _, err = db.AddRelicObtain(ctx, ormRelics); err != nil {
		return
	}
	return
}

// Convert a run stored in the database into a RunSchemaJson
func RunToJson(ctx context.Context, db *orm.Queries, play_id string) (data map[string]any, err error) {
	var rtr orm.RunToJsonRow
	if rtr, err = db.RunToJson(ctx, play_id); err != nil {
		return
	}
	data1 := make(map[string]any)
	extra := make(map[string]any)
	if err = rtr.RRaw.AssignTo(&data1); err != nil {
		return
	}
	if err = rtr.RExtra.AssignTo(&extra); err != nil {
		return
	}
	data1["path_per_floor"] = lo.Map(pathToStringRev(rtr.RPathPerFloor), func(v string, _ int) FloorPath {
		return ReNull(v)
	})
	data1["path_taken"] = pathToStringRev(rtr.RPathTaken)
	for k, v := range extra {
		data1[k] = v
	}
	return data1, nil
}

var pathToMapFwd = map[string]string{
	"BOSS": "B",
}
var pathToMapRev = lo.Invert(pathToMapFwd)

func pathToStringFwd(ar []string) string {
	out := make([]string, len(ar))
	for _, toS := range ar {
		var ch string
		if len(toS) == 1 {
			ch = toS
		} else if dstVal := pathToMapFwd[toS]; dstVal != "" {
			ch = dstVal
		} else {
			ch = fmt.Sprintf(",%s,", toS)
		}
		out = append(out, ch)
	}
	return strings.Join(out, "")
}

func pathToStringRev(s string) []string {
	out := make([]string, 0)
	parts := strings.Split(s, ",")
	for i, p := range parts {
		if i%2 == 0 {
			for _, r := range p {
				ch := string(r)
				if ch2 := pathToMapRev[ch]; ch2 != "" {
					ch = ch2
				}
				out = append(out, ch)
			}
		} else {
			out = append(out, p)
		}
	}
	return out
}

const NULL_STR_CHAR = "\x1B"

func DeNull(s *string) string {
	if s == nil {
		return NULL_STR_CHAR
	} else {
		return *s
	}
}

func ReNull(s string) *string {
	if s == NULL_STR_CHAR {
		return nil
	} else {
		return &s
	}
}

func sqlInt32(v int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(v), Valid: true}
}

// Convert an array of float64 to an array of int32
func mapInt32(ar []float64) []int32 {
	out := make([]int32, len(ar))
	for i, v := range ar {
		out[i] = int32(v)
	}
	return out
}
