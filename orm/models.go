// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package orm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
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
}

type Campfirechoice struct {
	ID    int32
	RunID int32
	Data  sql.NullInt32
	Floor int32
	Key   int32
}

type Cardchoice struct {
	ID        int32
	RunID     int32
	NotPicked []int32
	Picked    int32
	Floor     int32
}

type CharacterList struct {
	ID   int32
	Name string
}

type Damagetaken struct {
	ID      int32
	RunID   int32
	Enemies int32
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

type Masterdeck struct {
	ID       int32
	RunID    int32
	CardID   int32
	Upgrades int16
	Count    int16
}

type Perfloordatum struct {
	ID        int32
	RunID     int32
	Floor     int16
	Gold      int32
	CurrentHp int32
	MaxHp     int32
}

type Potionobtain struct {
	ID    int32
	RunID int32
	Floor int32
	Key   int32
}

type Relicobtain struct {
	ID    int32
	RunID int32
	Floor int32
	Key   int32
}

type Runflag struct {
	RunID int32
	Flag  FlagKind
}

type Runsdatum struct {
	ID                   int32
	AscensionLevel       int32
	BuildVersion         int32
	CampfireRested       sql.NullInt32
	CampfireUpgraded     sql.NullInt32
	CharacterChosen      int32
	ChooseSeed           bool
	CircletCount         sql.NullInt32
	FloorReached         int32
	Gold                 int32
	ItemsPurchasedFloors []int32
	ItemsPurchasedIds    []int32
	ItemsPurgedFloors    []int32
	ItemsPurgedIds       []int32
	KilledBy             int32
	LocalTime            string
	NeowBonusID          int32
	NeowCostID           int32
	PathPerFloor         string
	PathTaken            string
	PlayID               string
	PlayerExperience     int32
	Playtime             int32
	PotionsFloorSpawned  []int32
	PotionsFloorUsage    []int32
	PurchasedPurges      int32
	Score                int32
	SeedPlayed           string
	SeedSourceTimestamp  sql.NullInt32
	Timestamp            sql.NullTime
	Victory              bool
	WinRate              float64
}

type Scope struct {
	ID     int32
	Key    string
	Desc   string
	Parent sql.NullInt32
}

type StatsCardCount struct {
	CharID   int32
	CardID   int32
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

type User struct {
	ID    int32
	Email string
}

type UsersToScope struct {
	UserID  int32
	ScopeID int32
}
