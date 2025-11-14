package dbModels

import (
	"QuickStone/src/common"
	"QuickStone/src/storage/database"
	"time"
)

type Object struct {
	UserName   string `gorm:"not null;primaryKey"`
	BucketName string `gorm:"not null;primaryKey;"`
	Key        string `gorm:"not null;primaryKey"`
	ObjectType string
	Size       common.ObjectSizeT
	CreateTime time.Time `gorm:"not null;autoCreateTime"`

	Bucket Bucket `gorm:"foreignKey:UserName,BucketName;references:UserName,Name"`
}

func init() {
	if err := database.Client.AutoMigrate(&Object{}); err != nil {
		panic(err)
	}
}
