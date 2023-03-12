// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: db_add.sql

package orm

import (
	"context"
	"database/sql"

	"github.com/jackc/pgtype"
)

type AddBossRelicsParams struct {
	RunID     int32
	NotPicked []int32
	Picked    int32
	Ord       int16
}

type AddCampfireParams struct {
	RunID    int32
	StrData  sql.NullInt32
	CardData sql.NullInt32
	Floor    int32
	Key      int32
}

type AddCardChoiceParams struct {
	RunID     int32
	Floor     int32
	NotPicked []int32
	Picked    int32
}

type AddDamageTakenParams struct {
	RunID   int32
	Enemies int32
	Damage  float32
	Floor   int32
	Turns   int32
}

type AddEventChoicesParams struct {
	RunID             int32
	DamageDelta       int32
	EventNameID       int32
	Floor             int32
	GoldDelta         int32
	MaxHpDelta        int32
	PlayerChoiceID    int32
	RelicsObtainedIds []int32
}

const addFlag = `-- name: AddFlag :exec
INSERT INTO RunFlags (run_id, flag) VALUES ($1, $2)
`

type AddFlagParams struct {
	RunID int32
	Flag  FlagKind
}

func (q *Queries) AddFlag(ctx context.Context, arg AddFlagParams) error {
	_, err := q.db.Exec(ctx, addFlag, arg.RunID, arg.Flag)
	return err
}

type AddItemsPurchasedParams struct {
	RunID  int32
	CardID int32
	Floor  int16
}

type AddItemsPurgedParams struct {
	RunID  int32
	CardID int32
	Floor  int16
}

type AddMasterDeckParams struct {
	RunID  int32
	CardID int32
	Count  int16
}

type AddPerFloorParams struct {
	RunID     int32
	Floor     int16
	Gold      int32
	CurrentHp int32
	MaxHp     int32
}

type AddPotionObtainParams struct {
	RunID int32
	Floor int16
	Key   int32
}

type AddRelicObtainParams struct {
	RunID int32
	Floor int16
	Key   int32
}

type AddRunArraysParams struct {
	RunID               int32
	DailyMods           []int32
	PotionsFloorSpawned []int32
	PotionsFloorUsage   []int32
	RelicIds            []int32
}

const addRunRaw = `-- name: AddRunRaw :one
INSERT INTO RunsData
    (ascension_level, build_version, campfire_rested, campfire_upgraded, character_id, choose_seed,
     circlet_count, floor_reached, gold, killed_by, local_time, neow_bonus_id, neow_cost_id,
     path_per_floor, path_taken, play_id, player_experience, playtime, purchased_purges, score,
     seed_played, seed_source_timestamp, special_seed, "timestamp", victory, win_rate)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26
)
RETURNING RunsData.id
`

type AddRunRawParams struct {
	AscensionLevel      int32
	BuildVersion        int32
	CampfireRested      int32
	CampfireUpgraded    int32
	CharacterID         int32
	ChooseSeed          bool
	CircletCount        int32
	FloorReached        int32
	Gold                int32
	KilledBy            int32
	LocalTime           string
	NeowBonusID         int32
	NeowCostID          int32
	PathPerFloor        string
	PathTaken           string
	PlayID              string
	PlayerExperience    int32
	Playtime            int32
	PurchasedPurges     int32
	Score               int32
	SeedPlayed          string
	SeedSourceTimestamp sql.NullInt32
	SpecialSeed         int32
	Timestamp           sql.NullTime
	Victory             bool
	WinRate             float64
}

func (q *Queries) AddRunRaw(ctx context.Context, arg AddRunRawParams) (int32, error) {
	row := q.db.QueryRow(ctx, addRunRaw,
		arg.AscensionLevel,
		arg.BuildVersion,
		arg.CampfireRested,
		arg.CampfireUpgraded,
		arg.CharacterID,
		arg.ChooseSeed,
		arg.CircletCount,
		arg.FloorReached,
		arg.Gold,
		arg.KilledBy,
		arg.LocalTime,
		arg.NeowBonusID,
		arg.NeowCostID,
		arg.PathPerFloor,
		arg.PathTaken,
		arg.PlayID,
		arg.PlayerExperience,
		arg.Playtime,
		arg.PurchasedPurges,
		arg.Score,
		arg.SeedPlayed,
		arg.SeedSourceTimestamp,
		arg.SpecialSeed,
		arg.Timestamp,
		arg.Victory,
		arg.WinRate,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const addRunsExtra = `-- name: AddRunsExtra :exec
INSERT INTO runs_extra (run_id, extra) VALUES ($1,$2)
`

type AddRunsExtraParams struct {
	RunID int32
	Extra pgtype.JSONB
}

func (q *Queries) AddRunsExtra(ctx context.Context, arg AddRunsExtraParams) error {
	_, err := q.db.Exec(ctx, addRunsExtra, arg.RunID, arg.Extra)
	return err
}

const archiveAdd = `-- name: ArchiveAdd :exec
INSERT INTO RawJsonArchive(bdata, play_id) VALUES ($1, $2)
`

type ArchiveAddParams struct {
	Bdata  pgtype.JSON
	PlayID string
}

func (q *Queries) ArchiveAdd(ctx context.Context, arg ArchiveAddParams) error {
	_, err := q.db.Exec(ctx, archiveAdd, arg.Bdata, arg.PlayID)
	return err
}

const archiveBegin = `-- name: ArchiveBegin :many
UPDATE rawjsonarchive ra SET status = $1 WHERE status = 0 RETURNING ra.id, ra.bdata, ra.play_id, ra.status
`

func (q *Queries) ArchiveBegin(ctx context.Context, status int16) ([]Rawjsonarchive, error) {
	rows, err := q.db.Query(ctx, archiveBegin, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Rawjsonarchive
	for rows.Next() {
		var i Rawjsonarchive
		if err := rows.Scan(
			&i.ID,
			&i.Bdata,
			&i.PlayID,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const archiveComplete = `-- name: ArchiveComplete :many
UPDATE rawjsonarchive ra SET status = -1 WHERE status = $1 RETURNING ra.id
`

func (q *Queries) ArchiveComplete(ctx context.Context, status int16) ([]int32, error) {
	rows, err := q.db.Query(ctx, archiveComplete, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const doesRunExist = `-- name: DoesRunExist :one
SELECT count(id)::boolean FROM RunsData R WHERE R.play_id = $1
`

func (q *Queries) DoesRunExist(ctx context.Context, playID string) (bool, error) {
	row := q.db.QueryRow(ctx, doesRunExist, playID)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const getRun = `-- name: GetRun :one
SELECT id, ascension_level, build_version, campfire_rested, campfire_upgraded, character_id, choose_seed, circlet_count, floor_reached, gold, killed_by, local_time, neow_bonus_id, neow_cost_id, path_per_floor, path_taken, play_id, player_experience, playtime, purchased_purges, score, seed_played, seed_source_timestamp, special_seed, timestamp, victory, win_rate FROM RunsData WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRun(ctx context.Context, id int32) (Runsdatum, error) {
	row := q.db.QueryRow(ctx, getRun, id)
	var i Runsdatum
	err := row.Scan(
		&i.ID,
		&i.AscensionLevel,
		&i.BuildVersion,
		&i.CampfireRested,
		&i.CampfireUpgraded,
		&i.CharacterID,
		&i.ChooseSeed,
		&i.CircletCount,
		&i.FloorReached,
		&i.Gold,
		&i.KilledBy,
		&i.LocalTime,
		&i.NeowBonusID,
		&i.NeowCostID,
		&i.PathPerFloor,
		&i.PathTaken,
		&i.PlayID,
		&i.PlayerExperience,
		&i.Playtime,
		&i.PurchasedPurges,
		&i.Score,
		&i.SeedPlayed,
		&i.SeedSourceTimestamp,
		&i.SpecialSeed,
		&i.Timestamp,
		&i.Victory,
		&i.WinRate,
	)
	return i, err
}

const getStr = `-- name: GetStr :one
SELECT id FROM StrCache WHERE str = $1
`

func (q *Queries) GetStr(ctx context.Context, str string) (int32, error) {
	row := q.db.QueryRow(ctx, getStr, str)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const runToJson = `-- name: RunToJson :one
SELECT r.raw::json, r.path_per_floor::text, r.path_taken::text, r.extra::json
FROM run_to_json((SELECT id FROM runsdata WHERE play_id = $1)) r
`

type RunToJsonRow struct {
	RRaw          pgtype.JSON
	RPathPerFloor string
	RPathTaken    string
	RExtra        pgtype.JSON
}

func (q *Queries) RunToJson(ctx context.Context, playID string) (RunToJsonRow, error) {
	row := q.db.QueryRow(ctx, runToJson, playID)
	var i RunToJsonRow
	err := row.Scan(
		&i.RRaw,
		&i.RPathPerFloor,
		&i.RPathTaken,
		&i.RExtra,
	)
	return i, err
}

const strCacheAdd = `-- name: StrCacheAdd :exec
SELECT str_cache_add($1::text[])
`

func (q *Queries) StrCacheAdd(ctx context.Context, dollar_1 []string) error {
	_, err := q.db.Exec(ctx, strCacheAdd, dollar_1)
	return err
}

const strCacheToId = `-- name: StrCacheToId :many
SELECT str_cache_to_id($1::text[])
`

func (q *Queries) StrCacheToId(ctx context.Context, dollar_1 []string) ([]int32, error) {
	rows, err := q.db.Query(ctx, strCacheToId, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var str_cache_to_id int32
		if err := rows.Scan(&str_cache_to_id); err != nil {
			return nil, err
		}
		items = append(items, str_cache_to_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
