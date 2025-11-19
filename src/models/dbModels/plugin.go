package dbModels

import (
	"QuickStone/src/storage/database"
	"time"
)

type Plugin struct {
	Id         uint32 `gorm:"primaryKey;autoIncrement"`
	PluginName string `gorm:"not null;uniqueIndex"`
	UserName   string `gorm:"not null;"`
}

type PluginVersion struct {
	Id uint `gorm:"primaryKey"`

	PluginId uint32 `gorm:"uniqueIndex:uidx_plg_ver,priority:1;not null"`
	Plugin   Plugin `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// 版本规范化后的字符串，用来做唯一约束
	NormalizedVersion string `gorm:"size:120;not null;index:uniqueIndex:uidx_plg_ver,priority:2"`

	// wasm 二进制
	// Postgres: bytea
	Wasm []byte `gorm:"type:bytea"`

	Sha256    string `gorm:"size:64;index"`
	SizeBytes int64  `gorm:"not null;default:0"`

	CreatedAt time.Time `gorm:"index"`
}

func init() {
	if err := database.Client.AutoMigrate(&Plugin{}, &PluginVersion{}); err != nil {
		panic(err)
	}
}
