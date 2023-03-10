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

func (c BossRelicChoice) ToOrm(sc StrCache, runid int32) orm.AddBossRelicsParams {
	return orm.AddBossRelicsParams{
		RunID:     runid,
		NotPicked: sc.GetAll(c.NotPicked),
		Picked:    sc.Get(c.Picked),
	}
}

func (c *BossRelicChoice) GetStrings() []string {
	return append([]string{c.Picked}, c.NotPicked...)
}

// Convert CampfireChoice to an orm object
func (c CampfireChoice) ToOrm(sc StrCache, runid int32) orm.AddCampfireParams {
	return orm.AddCampfireParams{
		RunID: runid,
		Data:  sc.MaybeGet(c.Data),
		Key:   sc.Get(c.Key),
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

func (c CardChoice) ToOrm(sc StrCache, runid int32) orm.AddCardChoiceParams {
	return orm.AddCardChoiceParams{
		RunID:     runid,
		Floor:     int32(c.Floor),
		NotPicked: sc.GetAll(c.NotPicked),
		Picked:    sc.Get(c.Picked),
	}
}

func (c *CardChoice) GetStrings() []string {
	return append([]string{c.Picked}, c.NotPicked...)
}

func (c DamageTaken) ToOrm(sc StrCache, runid int32) orm.AddDamageTakenParams {
	return orm.AddDamageTakenParams{
		RunID:   runid,
		Floor:   int32(c.Floor),
		Turns:   int32(c.Turns),
		Damage:  float32(c.Damage),
		Enemies: sc.Get(c.Enemies),
	}
}

func (c *DamageTaken) GetStrings() []string {
	return []string{c.Enemies}
}

func (c EventChoice) ToOrm(sc StrCache, runid int32) orm.AddEventChoicesParams {
	return orm.AddEventChoicesParams{
		RunID:             runid,
		DamageDelta:       int32(c.DamageHealed - c.DamageTaken),
		EventNameID:       sc.Get(c.EventName),
		Floor:             int32(c.Floor),
		GoldDelta:         int32(c.GoldGain - c.GoldLoss),
		MaxHpDelta:        int32(c.MaxHpGain - c.MaxHpLoss),
		PlayerChoiceID:    sc.Get(c.PlayerChoice),
		RelicsObtainedIds: sc.GetAll(c.RelicsObtained),
	}
}

func (c *EventChoice) GetStrings() []string {
	return append([]string{c.EventName, c.PlayerChoice}, c.RelicsObtained...)
}

func (c PotionObtained) ToOrm(sc StrCache, runid int32) orm.AddPotionObtainParams {
	return orm.AddPotionObtainParams{
		RunID: runid,
		Floor: int16(c.Floor),
		Key:   sc.Get(c.Key),
	}
}

func (c *PotionObtained) GetStrings() []string {
	return []string{c.Key}
}

func (s RelicObtain) ToOrm(sc StrCache, runid int32) orm.AddRelicObtainParams {
	return orm.AddRelicObtainParams{
		RunID: runid,
		Floor: int16(s.Floor),
		Key:   sc.Get(s.Key),
	}
}

func (c *RelicObtain) GetStrings() []string {
	return []string{c.Key}
}

func (s *RunSchemaJson) ToAddRunRaw(sc StrCache) orm.AddRunRawParams {
	tstamp := time.UnixMilli(int64(s.Timestamp))
	pathNorm := lo.Map(s.PathPerFloor, func(v FloorPath, _ int) string {
		return DeNull(v)
	})

	return orm.AddRunRawParams{
		AscensionLevel:   int32(s.AscensionLevel),
		BuildVersion:     sc.Get(s.BuildVersion),
		CampfireRested:   int32(s.CampfireRested),
		CampfireUpgraded: int32(s.CampfireUpgraded),
		CharacterID:      sc.Get(s.CharacterChosen),
		ChooseSeed:       s.ChoseSeed,
		CircletCount:     int32(s.CircletCount),
		FloorReached:     int32(s.FloorReached),
		Gold:             int32(s.Gold),
		KilledBy:         sc.Get(s.KilledBy),
		LocalTime:        s.LocalTime,
		NeowBonusID:      sc.Get(s.NeowBonus),
		NeowCostID:       sc.Get(s.NeowCost),
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

func (s *RunSchemaJson) toPerFloorOrm(runid int32) []orm.AddPerFloorParams {
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

func (s *RunSchemaJson) toArrayOrm(sc StrCache, runid int32) []orm.AddRunArraysParams {
	return []orm.AddRunArraysParams{{
		RunID:                runid,
		DailyMods:            sc.GetAll(s.DailyMods),
		ItemsPurchasedFloors: mapInt32(s.ItemPurchaseFloors),
		ItemsPurchasedIds:    sc.GetAll(s.ItemsPurged),
		ItemsPurgedFloors:    mapInt32(s.ItemsPurgedFloors),
		ItemsPurgedIds:       sc.GetAll(s.ItemsPurged),
		PotionsFloorSpawned:  mapInt32(s.PotionsFloorSpawned),
		PotionsFloorUsage:    mapInt32(s.PotionsFloorUsage),
		RelicIds:             sc.GetAll(s.Relics),
	}}
}

func (s *RunSchemaJson) GetStrings() []string {
	out := []string{s.BuildVersion, s.CharacterChosen, s.KilledBy, s.NeowBonus, s.NeowCost}
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
	for _, u := range s.CardChoices {
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

// Add this Run to the database. Returns the rowid of the run.
func (r *RunSchemaJson) AddToDb(ctx context.Context, sc StrCache, db *orm.Queries) (runId int32, err error) {
	deck := NewMasterDeck()
	deck.AddCards(r.MasterDeck)
	if err = sc.Load(ctx, r.GetStrings(), deck.GetStrings()); err != nil {
		return
	}
	runId, err = db.AddRunRaw(ctx, r.ToAddRunRaw(sc))
	if err != nil {
		return
	}
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
	if _, err = db.AddRunArrays(ctx, r.toArrayOrm(sc, runId)); err != nil {
		return
	}
	if _, err = db.AddMasterDeck(ctx, deck.ToOrm(sc, runId)); err != nil {
		return
	}
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

	perFloor := r.toPerFloorOrm(runId)
	if _, err = db.AddPerFloor(ctx, perFloor); err != nil {
		return
	}
	ormBossRelics := MapToOrm[orm.AddBossRelicsParams](r.BossRelics, sc, runId)
	if _, err = db.AddBossRelics(ctx, ormBossRelics); err != nil {
		return
	}
	ormCampfires := MapToOrm[orm.AddCampfireParams](r.CampfireChoices, sc, runId)
	if _, err = db.AddCampfire(ctx, ormCampfires); err != nil {
		return
	}
	ormCards := MapToOrm[orm.AddCardChoiceParams](r.CardChoices, sc, runId)
	if _, err = db.AddCardChoice(ctx, ormCards); err != nil {
		return
	}
	ormDamageTaken := MapToOrm[orm.AddDamageTakenParams](r.DamageTaken, sc, runId)
	if _, err = db.AddDamageTaken(ctx, ormDamageTaken); err != nil {
		return
	}
	ormEvents := MapToOrm[orm.AddEventChoicesParams](r.EventChoices, sc, runId)
	if _, err = db.AddEventChoices(ctx, ormEvents); err != nil {
		return
	}
	ormPotions := MapToOrm[orm.AddPotionObtainParams](r.PotionsObtained, sc, runId)
	if _, err = db.AddPotionObtain(ctx, ormPotions); err != nil {
		return
	}
	ormRelics := MapToOrm[orm.AddRelicObtainParams](r.RelicsObtained, sc, runId)
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

type ConvToOrm[T any] interface {
	ToOrm(sc StrCache, runid int32) T
}

// Calls ToOrm on inp, returning the result
func MapToOrm[T any, E ConvToOrm[T]](inp []E, sc StrCache, runid int32) []T {
	out := make([]T, len(inp))
	for i, e := range inp {
		out[i] = e.ToOrm(sc, runid)
	}
	return out
}
