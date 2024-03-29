// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: copyfrom.go

package orm

import (
	"context"
)

// iteratorForAddBossRelics implements pgx.CopyFromSource.
type iteratorForAddBossRelics struct {
	rows                 []AddBossRelicsParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddBossRelics) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddBossRelics) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].NotPicked,
		r.rows[0].Picked,
		r.rows[0].Ord,
	}, nil
}

func (r iteratorForAddBossRelics) Err() error {
	return nil
}

func (q *Queries) AddBossRelics(ctx context.Context, arg []AddBossRelicsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"bossrelics"}, []string{"run_id", "not_picked", "picked", "ord"}, &iteratorForAddBossRelics{rows: arg})
}

// iteratorForAddCampfire implements pgx.CopyFromSource.
type iteratorForAddCampfire struct {
	rows                 []AddCampfireParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddCampfire) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddCampfire) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].StrData,
		r.rows[0].CardData,
		r.rows[0].Floor,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForAddCampfire) Err() error {
	return nil
}

func (q *Queries) AddCampfire(ctx context.Context, arg []AddCampfireParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"campfirechoice"}, []string{"run_id", "str_data", "card_data", "floor", "key"}, &iteratorForAddCampfire{rows: arg})
}

// iteratorForAddCardChoice implements pgx.CopyFromSource.
type iteratorForAddCardChoice struct {
	rows                 []AddCardChoiceParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddCardChoice) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddCardChoice) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].Floor,
		r.rows[0].NotPicked,
		r.rows[0].Picked,
	}, nil
}

func (r iteratorForAddCardChoice) Err() error {
	return nil
}

func (q *Queries) AddCardChoice(ctx context.Context, arg []AddCardChoiceParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"cardchoices"}, []string{"run_id", "floor", "not_picked", "picked"}, &iteratorForAddCardChoice{rows: arg})
}

// iteratorForAddDamageTaken implements pgx.CopyFromSource.
type iteratorForAddDamageTaken struct {
	rows                 []AddDamageTakenParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddDamageTaken) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddDamageTaken) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].Enemies,
		r.rows[0].Damage,
		r.rows[0].Floor,
		r.rows[0].Turns,
	}, nil
}

func (r iteratorForAddDamageTaken) Err() error {
	return nil
}

func (q *Queries) AddDamageTaken(ctx context.Context, arg []AddDamageTakenParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"damagetaken"}, []string{"run_id", "enemies", "damage", "floor", "turns"}, &iteratorForAddDamageTaken{rows: arg})
}

// iteratorForAddEventChoices implements pgx.CopyFromSource.
type iteratorForAddEventChoices struct {
	rows                 []AddEventChoicesParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddEventChoices) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddEventChoices) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].DamageDelta,
		r.rows[0].EventNameID,
		r.rows[0].Floor,
		r.rows[0].GoldDelta,
		r.rows[0].MaxHpDelta,
		r.rows[0].PlayerChoiceID,
		r.rows[0].RelicsObtainedIds,
	}, nil
}

func (r iteratorForAddEventChoices) Err() error {
	return nil
}

func (q *Queries) AddEventChoices(ctx context.Context, arg []AddEventChoicesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"eventchoices"}, []string{"run_id", "damage_delta", "event_name_id", "floor", "gold_delta", "max_hp_delta", "player_choice_id", "relics_obtained_ids"}, &iteratorForAddEventChoices{rows: arg})
}

// iteratorForAddItemsPurchased implements pgx.CopyFromSource.
type iteratorForAddItemsPurchased struct {
	rows                 []AddItemsPurchasedParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddItemsPurchased) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddItemsPurchased) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].CardID,
		r.rows[0].Floor,
	}, nil
}

func (r iteratorForAddItemsPurchased) Err() error {
	return nil
}

func (q *Queries) AddItemsPurchased(ctx context.Context, arg []AddItemsPurchasedParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"itemspurchased"}, []string{"run_id", "card_id", "floor"}, &iteratorForAddItemsPurchased{rows: arg})
}

// iteratorForAddItemsPurged implements pgx.CopyFromSource.
type iteratorForAddItemsPurged struct {
	rows                 []AddItemsPurgedParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddItemsPurged) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddItemsPurged) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].CardID,
		r.rows[0].Floor,
	}, nil
}

func (r iteratorForAddItemsPurged) Err() error {
	return nil
}

func (q *Queries) AddItemsPurged(ctx context.Context, arg []AddItemsPurgedParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"itemspurged"}, []string{"run_id", "card_id", "floor"}, &iteratorForAddItemsPurged{rows: arg})
}

// iteratorForAddPerFloor implements pgx.CopyFromSource.
type iteratorForAddPerFloor struct {
	rows                 []AddPerFloorParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddPerFloor) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddPerFloor) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].Floor,
		r.rows[0].Gold,
		r.rows[0].CurrentHp,
		r.rows[0].MaxHp,
	}, nil
}

func (r iteratorForAddPerFloor) Err() error {
	return nil
}

func (q *Queries) AddPerFloor(ctx context.Context, arg []AddPerFloorParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"perfloordata"}, []string{"run_id", "floor", "gold", "current_hp", "max_hp"}, &iteratorForAddPerFloor{rows: arg})
}

// iteratorForAddPotionObtain implements pgx.CopyFromSource.
type iteratorForAddPotionObtain struct {
	rows                 []AddPotionObtainParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddPotionObtain) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddPotionObtain) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].Floor,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForAddPotionObtain) Err() error {
	return nil
}

func (q *Queries) AddPotionObtain(ctx context.Context, arg []AddPotionObtainParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"potionobtains"}, []string{"run_id", "floor", "key"}, &iteratorForAddPotionObtain{rows: arg})
}

// iteratorForAddRelicObtain implements pgx.CopyFromSource.
type iteratorForAddRelicObtain struct {
	rows                 []AddRelicObtainParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddRelicObtain) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddRelicObtain) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].Floor,
		r.rows[0].Key,
	}, nil
}

func (r iteratorForAddRelicObtain) Err() error {
	return nil
}

func (q *Queries) AddRelicObtain(ctx context.Context, arg []AddRelicObtainParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"relicobtains"}, []string{"run_id", "floor", "key"}, &iteratorForAddRelicObtain{rows: arg})
}

// iteratorForAddRunArrays implements pgx.CopyFromSource.
type iteratorForAddRunArrays struct {
	rows                 []AddRunArraysParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddRunArrays) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddRunArrays) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RunID,
		r.rows[0].DailyMods,
		r.rows[0].MasterDeck,
		r.rows[0].PotionsFloorSpawned,
		r.rows[0].PotionsFloorUsage,
		r.rows[0].RelicIds,
	}, nil
}

func (r iteratorForAddRunArrays) Err() error {
	return nil
}

func (q *Queries) AddRunArrays(ctx context.Context, arg []AddRunArraysParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"runarrays"}, []string{"run_id", "daily_mods", "master_deck", "potions_floor_spawned", "potions_floor_usage", "relic_ids"}, &iteratorForAddRunArrays{rows: arg})
}
