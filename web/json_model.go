// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package web

import "encoding/json"
import "fmt"
import "github.com/google/uuid"
import "github.com/samber/lo"

type BossRelicChoice struct {
	// NotPicked corresponds to the JSON schema field "not_picked".
	NotPicked []string `json:"not_picked" yaml:"not_picked"`

	// Picked corresponds to the JSON schema field "picked".
	Picked string `json:"picked" yaml:"picked"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BossRelicChoice) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain BossRelicChoice
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["not_picked"]; !ok || v == nil {
		plain.NotPicked = []string{}
	}
	if v, ok := raw["picked"]; !ok || v == nil {
		plain.Picked = ""
	}
	*j = BossRelicChoice(plain)
	return nil
}

// One campfire selection
type CampfireChoice struct {
	// If KEY is 'SMITH', this will be the card ID that was upgraded
	Data *string `json:"data,omitempty" yaml:"data,omitempty"`

	// Which floor the campfire was on
	Floor float64 `json:"floor" yaml:"floor"`

	// Key corresponds to the JSON schema field "key".
	Key string `json:"key" yaml:"key"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CampfireChoice) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in CampfireChoice: required")
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key in CampfireChoice: required")
	}
	type Plain CampfireChoice
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CampfireChoice(plain)
	return nil
}

// One card choice
type CardChoice struct {
	// Floor corresponds to the JSON schema field "floor".
	Floor float64 `json:"floor" yaml:"floor"`

	// Cards that were not picked
	NotPicked []string `json:"not_picked" yaml:"not_picked"`

	// Picked corresponds to the JSON schema field "picked".
	Picked string `json:"picked" yaml:"picked"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CardChoice) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in CardChoice: required")
	}
	if v, ok := raw["not_picked"]; !ok || v == nil {
		return fmt.Errorf("field not_picked in CardChoice: required")
	}
	if v, ok := raw["picked"]; !ok || v == nil {
		return fmt.Errorf("field picked in CardChoice: required")
	}
	type Plain CardChoice
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CardChoice(plain)
	return nil
}

type DamageTaken struct {
	// How much damage the player took
	Damage float64 `json:"damage" yaml:"damage"`

	// Enemies corresponds to the JSON schema field "enemies".
	Enemies string `json:"enemies" yaml:"enemies"`

	// Which floor the fight occured on
	Floor float64 `json:"floor" yaml:"floor"`

	// How many turns the fight lasted
	Turns float64 `json:"turns" yaml:"turns"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DamageTaken) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["damage"]; !ok || v == nil {
		return fmt.Errorf("field damage in DamageTaken: required")
	}
	if v, ok := raw["enemies"]; !ok || v == nil {
		return fmt.Errorf("field enemies in DamageTaken: required")
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in DamageTaken: required")
	}
	if v, ok := raw["turns"]; !ok || v == nil {
		return fmt.Errorf("field turns in DamageTaken: required")
	}
	type Plain DamageTaken
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = DamageTaken(plain)
	return nil
}

type EventChoice struct {
	// HP gained
	DamageHealed float64 `json:"damage_healed" yaml:"damage_healed"`

	// DamageTaken corresponds to the JSON schema field "damage_taken".
	DamageTaken float64 `json:"damage_taken" yaml:"damage_taken"`

	// EventName corresponds to the JSON schema field "event_name".
	EventName string `json:"event_name" yaml:"event_name"`

	// Floor event occured on
	Floor float64 `json:"floor" yaml:"floor"`

	// GoldGain corresponds to the JSON schema field "gold_gain".
	GoldGain float64 `json:"gold_gain" yaml:"gold_gain"`

	// GoldLoss corresponds to the JSON schema field "gold_loss".
	GoldLoss float64 `json:"gold_loss" yaml:"gold_loss"`

	// MaxHpGain corresponds to the JSON schema field "max_hp_gain".
	MaxHpGain float64 `json:"max_hp_gain" yaml:"max_hp_gain"`

	// MaxHpLoss corresponds to the JSON schema field "max_hp_loss".
	MaxHpLoss float64 `json:"max_hp_loss" yaml:"max_hp_loss"`

	// PlayerChoice corresponds to the JSON schema field "player_choice".
	PlayerChoice string `json:"player_choice" yaml:"player_choice"`

	// RelicsObtained corresponds to the JSON schema field "relics_obtained".
	RelicsObtained []string `json:"relics_obtained" yaml:"relics_obtained"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EventChoice) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["damage_healed"]; !ok || v == nil {
		return fmt.Errorf("field damage_healed in EventChoice: required")
	}
	if v, ok := raw["damage_taken"]; !ok || v == nil {
		return fmt.Errorf("field damage_taken in EventChoice: required")
	}
	if v, ok := raw["event_name"]; !ok || v == nil {
		return fmt.Errorf("field event_name in EventChoice: required")
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in EventChoice: required")
	}
	if v, ok := raw["gold_gain"]; !ok || v == nil {
		return fmt.Errorf("field gold_gain in EventChoice: required")
	}
	if v, ok := raw["gold_loss"]; !ok || v == nil {
		return fmt.Errorf("field gold_loss in EventChoice: required")
	}
	if v, ok := raw["max_hp_gain"]; !ok || v == nil {
		return fmt.Errorf("field max_hp_gain in EventChoice: required")
	}
	if v, ok := raw["max_hp_loss"]; !ok || v == nil {
		return fmt.Errorf("field max_hp_loss in EventChoice: required")
	}
	if v, ok := raw["player_choice"]; !ok || v == nil {
		return fmt.Errorf("field player_choice in EventChoice: required")
	}
	type Plain EventChoice
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["relics_obtained"]; !ok || v == nil {
		plain.RelicsObtained = []string{}
	}
	*j = EventChoice(plain)
	return nil
}

type FloorPath *string

// When a potion was obtained
type PotionObtained struct {
	// Floor corresponds to the JSON schema field "floor".
	Floor float64 `json:"floor" yaml:"floor"`

	// Potion ID
	Key string `json:"key" yaml:"key"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PotionObtained) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in PotionObtained: required")
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key in PotionObtained: required")
	}
	type Plain PotionObtained
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PotionObtained(plain)
	return nil
}

type RelicObtain struct {
	// Floor corresponds to the JSON schema field "floor".
	Floor float64 `json:"floor" yaml:"floor"`

	// Key corresponds to the JSON schema field "key".
	Key string `json:"key" yaml:"key"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelicObtain) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["floor"]; !ok || v == nil {
		return fmt.Errorf("field floor in RelicObtain: required")
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key in RelicObtain: required")
	}
	type Plain RelicObtain
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = RelicObtain(plain)
	return nil
}

type RunSchemaJson struct {
	// The Ascension level (0 - 20) TODO - does this appear on 0 ascension?
	AscensionLevel int `json:"ascension_level" yaml:"ascension_level"`

	// BossRelics corresponds to the JSON schema field "boss_relics".
	BossRelics []BossRelicChoice `json:"boss_relics" yaml:"boss_relics"`

	// BuildVersion corresponds to the JSON schema field "build_version".
	BuildVersion string `json:"build_version" yaml:"build_version"`

	// CampfireChoices corresponds to the JSON schema field "campfire_choices".
	CampfireChoices []CampfireChoice `json:"campfire_choices" yaml:"campfire_choices"`

	// Number of times rested at a campfire
	CampfireRested float64 `json:"campfire_rested" yaml:"campfire_rested"`

	// CampfireUpgraded corresponds to the JSON schema field "campfire_upgraded".
	CampfireUpgraded float64 `json:"campfire_upgraded" yaml:"campfire_upgraded"`

	// List of card choices
	CardChoices []CardChoice `json:"card_choices" yaml:"card_choices"`

	// CharacterChosen corresponds to the JSON schema field "character_chosen".
	CharacterChosen string `json:"character_chosen" yaml:"character_chosen"`

	// Was the entered manually (true) or random (false)
	ChoseSeed bool `json:"chose_seed" yaml:"chose_seed"`

	// How many unknown relics the player had
	CircletCount int `json:"circlet_count" yaml:"circlet_count"`

	// CurrentHpPerFloor corresponds to the JSON schema field "current_hp_per_floor".
	CurrentHpPerFloor []float64 `json:"current_hp_per_floor" yaml:"current_hp_per_floor"`

	// List of modifiers for daily runs, only found if is_daily is true
	DailyMods []string `json:"daily_mods,omitempty" yaml:"daily_mods,omitempty"`

	// List of DamageTaken events
	DamageTaken []DamageTaken `json:"damage_taken" yaml:"damage_taken"`

	// EventChoices corresponds to the JSON schema field "event_choices".
	EventChoices []EventChoice `json:"event_choices" yaml:"event_choices"`

	// What floor number the player reached
	FloorReached float64 `json:"floor_reached" yaml:"floor_reached"`

	// Gold amount at the end of the run
	Gold float64 `json:"gold" yaml:"gold"`

	// How much gold was obtained on each floor
	GoldPerFloor []float64 `json:"gold_per_floor" yaml:"gold_per_floor"`

	// Is the player playing on some ascension mode?
	IsAscensionMode bool `json:"is_ascension_mode" yaml:"is_ascension_mode"`

	// IsBeta corresponds to the JSON schema field "is_beta".
	IsBeta bool `json:"is_beta" yaml:"is_beta"`

	// IsDaily corresponds to the JSON schema field "is_daily".
	IsDaily bool `json:"is_daily,omitempty" yaml:"is_daily,omitempty"`

	// Is Endless mode
	IsEndless bool `json:"is_endless" yaml:"is_endless"`

	// IsProd corresponds to the JSON schema field "is_prod".
	IsProd bool `json:"is_prod,omitempty" yaml:"is_prod,omitempty"`

	// Is this a custom trial run
	IsTrial bool `json:"is_trial" yaml:"is_trial"`

	// ItemPurchaseFloors corresponds to the JSON schema field "item_purchase_floors".
	ItemPurchaseFloors []float64 `json:"item_purchase_floors" yaml:"item_purchase_floors"`

	// ItemsPurchased corresponds to the JSON schema field "items_purchased".
	ItemsPurchased []string `json:"items_purchased" yaml:"items_purchased"`

	// List of removed card IDs
	ItemsPurged []string `json:"items_purged" yaml:"items_purged"`

	// ItemsPurgedFloors corresponds to the JSON schema field "items_purged_floors".
	ItemsPurgedFloors []float64 `json:"items_purged_floors" yaml:"items_purged_floors"`

	// Encounter ID where player died
	KilledBy string `json:"killed_by" yaml:"killed_by"`

	// Local time in YYYYmmddHHMMSS format
	LocalTime string `json:"local_time" yaml:"local_time"`

	// List of card IDs in the master deck, +1 indicates an upgraded card
	MasterDeck []string `json:"master_deck" yaml:"master_deck"`

	// MaxHpPerFloor corresponds to the JSON schema field "max_hp_per_floor".
	MaxHpPerFloor []float64 `json:"max_hp_per_floor" yaml:"max_hp_per_floor"`

	// ID of player's Neow choice
	NeowBonus string `json:"neow_bonus" yaml:"neow_bonus"`

	// TODO
	NeowCost string `json:"neow_cost" yaml:"neow_cost"`

	// PathPerFloor corresponds to the JSON schema field "path_per_floor".
	PathPerFloor []FloorPath `json:"path_per_floor" yaml:"path_per_floor"`

	// Path the player took
	PathTaken []string `json:"path_taken" yaml:"path_taken"`

	// UUID for this run
	PlayId uuid.UUID `json:"play_id" yaml:"play_id"`

	// XP gained at the end of the run
	PlayerExperience float64 `json:"player_experience" yaml:"player_experience"`

	// Play time in seconds
	Playtime float64 `json:"playtime" yaml:"playtime"`

	// PotionsFloorSpawned corresponds to the JSON schema field
	// "potions_floor_spawned".
	PotionsFloorSpawned []float64 `json:"potions_floor_spawned" yaml:"potions_floor_spawned"`

	// Which floors the player used a potion on
	PotionsFloorUsage []float64 `json:"potions_floor_usage" yaml:"potions_floor_usage"`

	// PotionsObtained corresponds to the JSON schema field "potions_obtained".
	PotionsObtained []PotionObtained `json:"potions_obtained" yaml:"potions_obtained"`

	// PurchasedPurges corresponds to the JSON schema field "purchased_purges".
	PurchasedPurges int `json:"purchased_purges" yaml:"purchased_purges"`

	// Relics corresponds to the JSON schema field "relics".
	Relics []string `json:"relics" yaml:"relics"`

	// RelicsObtained corresponds to the JSON schema field "relics_obtained".
	RelicsObtained []RelicObtain `json:"relics_obtained" yaml:"relics_obtained"`

	// Player's score at the end of the run
	Score float64 `json:"score" yaml:"score"`

	// The run seed
	SeedPlayed string `json:"seed_played" yaml:"seed_played"`

	// SeedSourceTimestamp corresponds to the JSON schema field
	// "seed_source_timestamp".
	SeedSourceTimestamp int `json:"seed_source_timestamp" yaml:"seed_source_timestamp"`

	// bitwise OR of special run modifiers
	SpecialSeed float64 `json:"special_seed,omitempty" yaml:"special_seed,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp int `json:"timestamp" yaml:"timestamp"`

	// Victory corresponds to the JSON schema field "victory".
	Victory bool `json:"victory" yaml:"victory"`

	// WinRate corresponds to the JSON schema field "win_rate".
	WinRate float64 `json:"win_rate" yaml:"win_rate"`
	// Additional fields
	Extra map[string]any `json:"-"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RunSchemaJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["character_chosen"]; !ok || v == nil {
		return fmt.Errorf("field character_chosen in RunSchemaJson: required")
	}
	if v, ok := raw["gold_per_floor"]; !ok || v == nil {
		return fmt.Errorf("field gold_per_floor in RunSchemaJson: required")
	}
	if v, ok := raw["local_time"]; !ok || v == nil {
		return fmt.Errorf("field local_time in RunSchemaJson: required")
	}
	if v, ok := raw["max_hp_per_floor"]; !ok || v == nil {
		return fmt.Errorf("field max_hp_per_floor in RunSchemaJson: required")
	}
	if v, ok := raw["neow_cost"]; !ok || v == nil {
		return fmt.Errorf("field neow_cost in RunSchemaJson: required")
	}
	if v, ok := raw["path_per_floor"]; !ok || v == nil {
		return fmt.Errorf("field path_per_floor in RunSchemaJson: required")
	}
	if v, ok := raw["play_id"]; !ok || v == nil {
		return fmt.Errorf("field play_id in RunSchemaJson: required")
	}
	if v, ok := raw["seed_played"]; !ok || v == nil {
		return fmt.Errorf("field seed_played in RunSchemaJson: required")
	}
	type Plain RunSchemaJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["ascension_level"]; !ok || v == nil {
		plain.AscensionLevel = 0
	}
	if v, ok := raw["boss_relics"]; !ok || v == nil {
		plain.BossRelics = []BossRelicChoice{}
	}
	if v, ok := raw["build_version"]; !ok || v == nil {
		plain.BuildVersion = ""
	}
	if v, ok := raw["campfire_choices"]; !ok || v == nil {
		plain.CampfireChoices = []CampfireChoice{}
	}
	if v, ok := raw["campfire_rested"]; !ok || v == nil {
		plain.CampfireRested = 0
	}
	if v, ok := raw["campfire_upgraded"]; !ok || v == nil {
		plain.CampfireUpgraded = 0
	}
	if v, ok := raw["card_choices"]; !ok || v == nil {
		plain.CardChoices = []CardChoice{}
	}
	if v, ok := raw["chose_seed"]; !ok || v == nil {
		plain.ChoseSeed = false
	}
	if v, ok := raw["circlet_count"]; !ok || v == nil {
		plain.CircletCount = 0
	}
	if v, ok := raw["current_hp_per_floor"]; !ok || v == nil {
		plain.CurrentHpPerFloor = []float64{}
	}
	if v, ok := raw["damage_taken"]; !ok || v == nil {
		plain.DamageTaken = []DamageTaken{}
	}
	if v, ok := raw["event_choices"]; !ok || v == nil {
		plain.EventChoices = []EventChoice{}
	}
	if v, ok := raw["floor_reached"]; !ok || v == nil {
		plain.FloorReached = 0
	}
	if v, ok := raw["gold"]; !ok || v == nil {
		plain.Gold = 0
	}
	if v, ok := raw["is_ascension_mode"]; !ok || v == nil {
		plain.IsAscensionMode = false
	}
	if v, ok := raw["is_beta"]; !ok || v == nil {
		plain.IsBeta = false
	}
	if v, ok := raw["is_daily"]; !ok || v == nil {
		plain.IsDaily = false
	}
	if v, ok := raw["is_endless"]; !ok || v == nil {
		plain.IsEndless = false
	}
	if v, ok := raw["is_prod"]; !ok || v == nil {
		plain.IsProd = false
	}
	if v, ok := raw["is_trial"]; !ok || v == nil {
		plain.IsTrial = false
	}
	if v, ok := raw["item_purchase_floors"]; !ok || v == nil {
		plain.ItemPurchaseFloors = []float64{}
	}
	if v, ok := raw["items_purchased"]; !ok || v == nil {
		plain.ItemsPurchased = []string{}
	}
	if v, ok := raw["items_purged"]; !ok || v == nil {
		plain.ItemsPurged = []string{}
	}
	if v, ok := raw["items_purged_floors"]; !ok || v == nil {
		plain.ItemsPurgedFloors = []float64{}
	}
	if v, ok := raw["killed_by"]; !ok || v == nil {
		plain.KilledBy = ""
	}
	if v, ok := raw["master_deck"]; !ok || v == nil {
		plain.MasterDeck = []string{}
	}
	if v, ok := raw["neow_bonus"]; !ok || v == nil {
		plain.NeowBonus = ""
	}
	if v, ok := raw["path_taken"]; !ok || v == nil {
		plain.PathTaken = []string{}
	}
	if v, ok := raw["player_experience"]; !ok || v == nil {
		plain.PlayerExperience = 0
	}
	if v, ok := raw["playtime"]; !ok || v == nil {
		plain.Playtime = 0
	}
	if v, ok := raw["potions_floor_spawned"]; !ok || v == nil {
		plain.PotionsFloorSpawned = []float64{}
	}
	if v, ok := raw["potions_floor_usage"]; !ok || v == nil {
		plain.PotionsFloorUsage = []float64{}
	}
	if v, ok := raw["potions_obtained"]; !ok || v == nil {
		plain.PotionsObtained = []PotionObtained{}
	}
	if v, ok := raw["purchased_purges"]; !ok || v == nil {
		plain.PurchasedPurges = 0
	}
	if v, ok := raw["relics"]; !ok || v == nil {
		plain.Relics = []string{}
	}
	if v, ok := raw["relics_obtained"]; !ok || v == nil {
		plain.RelicsObtained = []RelicObtain{}
	}
	if v, ok := raw["score"]; !ok || v == nil {
		plain.Score = 0
	}
	if v, ok := raw["seed_source_timestamp"]; !ok || v == nil {
		plain.SeedSourceTimestamp = 0
	}
	if v, ok := raw["special_seed"]; !ok || v == nil {
		plain.SpecialSeed = 0
	}
	if v, ok := raw["timestamp"]; !ok || v == nil {
		plain.Timestamp = 0
	}
	if v, ok := raw["victory"]; !ok || v == nil {
		plain.Victory = false
	}
	if v, ok := raw["win_rate"]; !ok || v == nil {
		plain.WinRate = 0
	}
	plain.Extra = lo.OmitByKeys(raw, runSchemaJsonKeys)
	*j = RunSchemaJson(plain)
	return nil
}

var runSchemaJsonKeys = []string{
	"ascension_level",
	"boss_relics",
	"build_version",
	"campfire_choices",
	"campfire_rested",
	"campfire_upgraded",
	"card_choices",
	"character_chosen",
	"chose_seed",
	"circlet_count",
	"current_hp_per_floor",
	"daily_mods",
	"damage_taken",
	"event_choices",
	"floor_reached",
	"gold",
	"gold_per_floor",
	"is_ascension_mode",
	"is_beta",
	"is_daily",
	"is_endless",
	"is_prod",
	"is_trial",
	"item_purchase_floors",
	"items_purchased",
	"items_purged",
	"items_purged_floors",
	"killed_by",
	"local_time",
	"master_deck",
	"max_hp_per_floor",
	"neow_bonus",
	"neow_cost",
	"path_per_floor",
	"path_taken",
	"play_id",
	"player_experience",
	"playtime",
	"potions_floor_spawned",
	"potions_floor_usage",
	"potions_obtained",
	"purchased_purges",
	"relics",
	"relics_obtained",
	"score",
	"seed_played",
	"seed_source_timestamp",
	"special_seed",
	"timestamp",
	"victory",
	"win_rate",
}