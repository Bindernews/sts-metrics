package stms

// Helpers for mapping the generated JSON structs to the generated SQL structs

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bindernews/sts-msr/orm"
	"github.com/samber/lo"
)

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
	return []string{c.Key}
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
		Floor: int32(c.Floor),
		Key:   sc.Get(c.Key),
	}
}

func (c *PotionObtained) GetStrings() []string {
	return []string{c.Key}
}

func (s RelicObtain) ToOrm(sc StrCache, runid int32) orm.AddRelicObtainParams {
	return orm.AddRelicObtainParams{
		RunID: runid,
		Floor: int32(s.Floor),
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
		CampfireRested:   SqlInt32(int(s.CampfireRested)),
		CampfireUpgraded: SqlInt32(int(s.CampfireUpgraded)),
		ChooseSeed:       s.ChoseSeed,
		CircletCount:     SqlInt32(s.CircletCount),
		FloorReached:     int32(s.FloorReached),
		Gold:             int32(s.Gold),
		//
		ItemsPurchasedFloors: mapInt32(s.ItemPurchaseFloors),
		ItemsPurgedFloors:    mapInt32(s.ItemsPurgedFloors),
		//
		LocalTime:        s.LocalTime,
		PathPerFloor:     pathToStringFwd(pathNorm),
		PathTaken:        pathToStringFwd(s.PathTaken),
		PlayID:           s.PlayId.String(),
		PlayerExperience: int32(s.PlayerExperience),
		Playtime:         int32(s.Playtime),
		//
		PotionsFloorSpawned: mapInt32(s.PotionsFloorSpawned),
		PotionsFloorUsage:   mapInt32(s.PotionsFloorUsage),
		PurchasedPurges:     int32(s.PurchasedPurges),
		Score:               int32(s.Score),
		SeedPlayed:          s.SeedPlayed,
		SeedSourceTimestamp: SqlInt32(s.SeedSourceTimestamp),
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

func (s *RunSchemaJson) GetStrings() []string {
	out := []string{s.BuildVersion, s.CharacterChosen, s.KilledBy, s.NeowBonus, s.NeowCost}
	out = append(out, s.ItemsPurchased...)
	out = append(out, s.ItemsPurged...)
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
	if err = db.AddRunText(ctx, orm.AddRunTextParams{
		ID:                runId,
		BuildVersion:      sc.Get(r.BuildVersion),
		CharacterChosen:   sc.Get(r.CharacterChosen),
		ItemsPurchasedIds: sc.GetAll(r.ItemsPurchased),
		ItemsPurgedIds:    sc.GetAll(r.ItemsPurged),
		KilledBy:          sc.Get(r.KilledBy),
		NeowBonusID:       sc.Get(r.NeowBonus),
		NeowCostID:        sc.Get(r.NeowCost),
	}); err != nil {
		return
	}
	if _, err = db.AddMasterDeck(ctx, deck.ToOrm(sc, runId)); err != nil {
		return
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

func (s *RunSchemaJson) FromAddRunRaw(m *orm.AddRunRawParams) error {
	s.PathPerFloor = lo.Map(pathToStringRev(m.PathPerFloor), func(s string, _ int) FloorPath {
		return FloorPath(ReNull(s))
	})
	s.PathTaken = pathToStringRev(m.PathTaken)
	// TODO update
	return nil
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

func SqlInt32(v int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(v), Valid: true}
}

func SqlString(v *string) sql.NullString {
	if v == nil || *v == "" {
		return sql.NullString{Valid: false}
	} else {
		return sql.NullString{Valid: true, String: *v}
	}
}

func mapInt32(ar []float64) []int32 {
	out := make([]int32, len(ar))
	for i, v := range ar {
		out[i] = int32(v)
	}
	return out
}
