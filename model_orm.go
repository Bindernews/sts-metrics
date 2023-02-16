package stms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bindernews/sts-msr/orm"
	"github.com/samber/lo"
)

func (s *CampfireChoice) ToOrm(runid int32) orm.AddCampfireParams {
	return orm.AddCampfireParams{
		RunID: runid,
		Cdata: SqlString(s.Data),
		Ckey:  s.Key,
		Floor: int32(s.Floor),
	}
}

func (s *CardChoice) ToOrm(runid int32) orm.AddCardChoiceParams {
	return orm.AddCardChoiceParams{
		RunID:     runid,
		Floor:     int32(s.Floor),
		NotPicked: s.NotPicked,
		Picked:    s.Picked,
	}
}

func (s *DamageTaken) ToOrm(runid int32) orm.AddDamageTakenParams {
	return orm.AddDamageTakenParams{
		RunID:   runid,
		Floor:   int32(s.Floor),
		Turns:   int32(s.Turns),
		Enemies: s.Enemies,
	}
}

func (s *PotionObtained) ToOrm(runid int32) orm.AddPotionObtainParams {
	return orm.AddPotionObtainParams{
		RunID: runid,
		Floor: int32(s.Floor),
		Ckey:  s.Key,
	}
}

func (s *RelicObtain) ToOrm(runid int32) orm.AddRelicObtainParams {
	return orm.AddRelicObtainParams{
		RunID: runid,
		Floor: int32(s.Floor),
		Ckey:  s.Key,
	}
}

func (s *RunSchemaJson) ToAddRunRaw() orm.AddRunRawParams {
	tstamp := time.UnixMilli(int64(s.Timestamp))
	pathNorm := lo.Map(s.PathPerFloor, func(v FloorPath, _ int) string {
		return DeNull(v)
	})

	return orm.AddRunRawParams{
		AscensionLevel:    int32(s.AscensionLevel),
		CampfireRested:    SqlInt32(s.CampfireRested),
		CampfireUpgraded:  SqlInt32(s.CampfireUpgraded),
		ChooseSeed:        s.ChoseSeed,
		CircletCount:      SqlInt32(s.CircletCount),
		CurrentHpPerFloor: mapInt32(s.CurrentHpPerFloor),
		FloorReached:      int32(s.FloorReached),
		Gold:              int32(s.Gold),
		GoldPerFloor:      mapInt32(s.GoldPerFloor),
		IsBeta:            s.IsBeta,
		IsDaily:           s.IsDaily,
		IsEndless:         s.IsEndless,
		IsProd:            s.IsProd,
		IsTrial:           s.IsTrial,
		//
		ItemsPurchasedFloors: mapInt32(s.ItemPurchaseFloors),
		ItemsPurgedFloors:    mapInt32(s.ItemsPurgedFloors),
		//
		LocalTime:        s.LocalTime,
		MaxHpPerFloor:    mapInt32(s.MaxHpPerFloor),
		NeowBonus:        s.NeowBonus,
		NeowCost:         s.NeowCost,
		PathPerFloor:     pathToStringFwd(pathNorm),
		PathTaken:        pathToStringFwd(s.PathTaken),
		PlayID:           s.PlayId,
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
		WinRate:             int32(s.WinRate),
	}
}

func (s *RunSchemaJson) ToSetRunText(runid int32) orm.SetRunTextParams {
	return orm.SetRunTextParams{
		ID:                  runid,
		BuildVersion:        s.BuildVersion,
		CharacterChosen:     s.CharacterChosen,
		ItemsPurchasedNames: s.ItemsPurchased,
		ItemsPurgedNames:    s.ItemsPurged,
	}
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

func (r *RunSchemaJson) AddToDb(ctx context.Context, db *orm.Queries) (runId int32, err error) {
	runId, err = db.AddRunRaw(ctx, r.ToAddRunRaw())
	if err != nil {
		return
	}
	if err = db.SetRunText(ctx, r.ToSetRunText(runId)); err != nil {
		return
	}
	for _, cc := range r.CampfireChoices {
		if _, err = db.AddCampfire(ctx, cc.ToOrm(runId)); err != nil {
			return
		}
	}
	for _, cc := range r.CardChoices {
		if _, err = db.AddCardChoice(ctx, cc.ToOrm(runId)); err != nil {
			return
		}
	}
	for _, cc := range r.DamageTaken {
		if _, err = db.AddDamageTaken(ctx, cc.ToOrm(runId)); err != nil {
			return
		}
	}
	// TODO event choices
	//for _, cc := range r.EventChoices {}
	for _, cc := range r.PotionsObtained {
		if err = db.AddPotionObtain(ctx, cc.ToOrm(runId)); err != nil {
			return
		}
	}
	for _, cc := range r.RelicsObtained {
		if err = db.AddRelicObtain(ctx, cc.ToOrm(runId)); err != nil {
			return
		}
	}
	return
}

func SqlInt32(v int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(v), Valid: true}
}

func SqlString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{Valid: false}
	} else {
		return sql.NullString{Valid: true, String: v}
	}
}

func mapInt32(ar []int) []int32 {
	out := make([]int32, len(ar))
	for i, v := range ar {
		out[i] = int32(v)
	}
	return out
}
