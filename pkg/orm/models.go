// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package orm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

type FlagKind string

const (
	FlagKindAscension FlagKind = "ascension"
	FlagKindBeta      FlagKind = "beta"
	FlagKindDaily     FlagKind = "daily"
	FlagKindEndless   FlagKind = "endless"
	FlagKindProd      FlagKind = "prod"
	FlagKindTrial     FlagKind = "trial"
)

func (e *FlagKind) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FlagKind(s)
	case string:
		*e = FlagKind(s)
	default:
		return fmt.Errorf("unsupported scan type for FlagKind: %T", src)
	}
	return nil
}

type NullFlagKind struct {
	FlagKind FlagKind
	Valid    bool // Valid is true if FlagKind is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFlagKind) Scan(value interface{}) error {
	if value == nil {
		ns.FlagKind, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FlagKind.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFlagKind) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FlagKind), nil
}

type Bossrelic struct {
	ID        int32
	RunID     int32
	NotPicked []int32
	Picked    int32
	Ord       int16
}

type Campfirechoice struct {
	ID       int32
	RunID    int32
	StrData  sql.NullInt32
	CardData sql.NullInt32
	Floor    int32
	Key      int32
}

type Campfirechoicesstring struct {
	ID    int32
	RunID int32
	Data  interface{}
	Floor int32
	Key   string
}

type Cardchoice struct {
	ID        int32
	RunID     int32
	NotPicked []int32
	Picked    int32
	Floor     int32
}

type Cardspec struct {
	ID       int32
	Card     string
	Upgrades int32
}

type Cardspecsex struct {
	ID       int32
	Card     string
	Upgrades int32
	Suffix   string
	CardFull interface{}
}

type Cardspecsnew struct {
	ID    int32
	Added time.Time
}

type CharacterList struct {
	ID   int32
	Name string
}

type Damagetaken struct {
	ID      int32
	RunID   int32
	Enemies int32
	Damage  float32
	Floor   int32
	Turns   int32
}

type Eventchoice struct {
	ID                int32
	RunID             int32
	DamageDelta       int32
	EventNameID       int32
	Floor             int32
	GoldDelta         int32
	MaxHpDelta        int32
	PlayerChoiceID    int32
	RelicsObtainedIds []int32
}

type Itemspurchased struct {
	RunID  int32
	CardID int32
	Floor  int16
}

type Itemspurged struct {
	RunID  int32
	CardID int32
	Floor  int16
}

type Perfloordatum struct {
	RunID     int32
	Floor     int16
	Gold      int32
	CurrentHp int32
	MaxHp     int32
}

type Potionobtain struct {
	ID    int32
	RunID int32
	Floor int16
	Key   int32
}

type Rawjsonarchive struct {
	ID     int32
	Bdata  pgtype.JSON
	PlayID string
	Status int16
}

type Relicobtain struct {
	ID    int32
	RunID int32
	Floor int16
	Key   int32
}

type Runarray struct {
	RunID               int32
	DailyMods           []int32
	MasterDeck          []int32
	PotionsFloorSpawned []int32
	PotionsFloorUsage   []int32
	RelicIds            []int32
}

type Runarraysext struct {
	RunID               int32
	DailyMods           []string
	MasterDeck          []string
	PotionsFloorSpawned []int32
	PotionsFloorUsage   []int32
	Relics              interface{}
}

type Runflag struct {
	RunID int32
	Flag  FlagKind
}

type RunsExtra struct {
	RunID int32
	Extra pgtype.JSONB
}

type Runsdatum struct {
	ID                  int32
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
	Added               time.Time
}

type Scope struct {
	ID     int32
	Key    string
	Desc   string
	Parent sql.NullInt32
}

type StatsCardCount struct {
	CharID   int32
	CardID   interface{}
	Total    int64
	Upgrades int64
}

type StatsOverview struct {
	ID            int32
	Name          string
	Runs          int64
	Wins          int64
	AvgWinRate    float64
	PDeckSize     []float32
	PFloorReached []float32
}

type Strcache struct {
	ID  int32
	Str string
}

type Strcachenew struct {
	ID    int32
	Added time.Time
}

type User struct {
	ID    int32
	Email string
}

type UsersToScope struct {
	UserID  int32
	ScopeID int32
}