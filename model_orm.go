package stms

import (
	"context"
	"database/sql"

	"github.com/bindernews/sts-msr/orm"
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

func (s *RunSchemaJson) ToOrm() orm.AddRunParams {
	return orm.AddRunParams{
		AscensionLevel: int32(s.AscensionLevel),
		S:              s.BuildVersion,
		CampfireRested: SqlInt32(s.CampfireRested),
		// TODO remaining fields
	}
}

func (r *RunSchemaJson) AddToDb(ctx context.Context, db *orm.Queries) (runId int32, err error) {
	runId, err = db.AddRun(ctx, r.ToOrm())
	if err != nil {
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
