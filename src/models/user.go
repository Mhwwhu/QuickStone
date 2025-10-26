package models

import (
	"QuickStone/src/common"
	"QuickStone/src/storage/database"
)

type User struct {
	Id       common.UserIdT `gorm:"not null;primarykey;autoIncrement"`
	Name     string         `gorm:"not null;unique;size:32;index"`
	Password string         `gorm:"not null;size:64"`
}

func init() {
	if err := database.Client.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
