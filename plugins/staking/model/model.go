package model

import (
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Payout struct {
	ID             uint                `gorm:"primary_key" json:"-"`
	Account        address.SS58Address `gorm:"index;type:char(48);default: null;size:100" json:"account"`
	Amount         decimal.Decimal     `gorm:"type:decimal(30,0);" json:"amount"`
	Era            uint32              `gorm:"index" json:"era"`
	Stash          address.SS58Address `gorm:"index;type:char(48);default: null;size:100" json:"stash"`
	ValidatorStash address.SS58Address `gorm:"index;type:char(48);default: null;size:100" json:"validator_stash"`
	BlockTimestamp uint64              `gorm:"index" json:"block_timestamp"`
	ModuleId       string              `gorm:"type:varchar(50)" json:"module_id"`
	EventId        string              `gorm:"type:varchar(50)" json:"event_id"`
	SlashKton      bool                `json:"slash_kton"`
	ExtrinsicIndex string              `gorm:"type:varchar(22)" json:"extrinsic_index"`
	EventIndex     string              `gorm:"type:varchar(22)" json:"event_index"`
	Claimed        bool                `gorm:"index" json:"-"`
}

type PoolPayout struct {
	ID             uint                `gorm:"primary_key" json:"-"`
	Account        address.SS58Address `gorm:"index;type:char(48);default: null;size:100" json:"account"`
	Amount         decimal.Decimal     `gorm:"type:decimal(30,0);" json:"amount"`
	PoolId         uint32              `gorm:"index" json:"pool_id"`
	ModuleId       string              `gorm:"type:varchar(50)" json:"module_id"`
	EventId        string              `gorm:"type:varchar(50)" json:"event_id"`
	ExtrinsicIndex string              `gorm:"type:varchar(22)" json:"extrinsic_index"`
	EventIndex     string              `gorm:"type:varchar(22)" json:"event_index"`
	BlockTimestamp uint64              `gorm:"index:,sort:desc" json:"block_timestamp"`
}

type ValidatorPrefs struct {
	ID                uint                `gorm:"primary_key" json:"-"`
	Account           address.SS58Address `gorm:"index;type:char(48);unique;default: null;size:100" json:"account"`
	Commission        decimal.Decimal     `sql:"type:decimal(12,11);" json:"commission"`
	BlockedNomination bool                `json:"blocked_nomination"`
	Era               uint32              `gorm:"index"`
}

type EraInfo struct {
	ID               uint            `gorm:"primary_key" json:"-"`
	Era              uint32          `gorm:"index" json:"era"`
	TotalStake       decimal.Decimal `gorm:"type:decimal(30,0);" json:"total_stake"`
	Stakes           datatypes.JSONSlice[EraStake]
	TotalPoints      uint32
	TotalRewards     decimal.Decimal
	ValidatorPoints  datatypes.JSONType[map[address.SS58Address]uint32]
	ValidatorRewards datatypes.JSONType[map[address.SS58Address]decimal.Decimal]
	StakerRewards    datatypes.JSONType[map[address.SS58Address]decimal.Decimal]
	StartBlock       uint `gorm:"index"`
	EndBlock         uint `gorm:"index"`
}

type EraStake struct {
	Validator      address.SS58Address
	Staker         address.SS58Address
	Amount         decimal.Decimal
	ValidatorTotal decimal.Decimal
}
