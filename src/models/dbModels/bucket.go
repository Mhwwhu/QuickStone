package dbModels

import (
	"QuickStone/src/storage/database"
	"time"
)

type StorageType int
type ACLType int
type storageTypeUtil struct{}
type aclTypeUtil struct{}

var StorageTypeUtil storageTypeUtil
var ACLTypeUtil aclTypeUtil

const (
	PRIVATE ACLType = iota
	PUBLIC_READ
	PUBLIC
)

const (
	STANDARD StorageType = iota
	LOW_FREQ
)

type Bucket struct {
	UserName    string      `gorm:"not null;primaryKey"`
	Name        string      `gorm:"not null;primaryKey;"`
	Area        string      `gorm:"not null;index"`
	StorageType StorageType `gorm:"not null"`
	ACLType     ACLType     `gorm:"not null"`
	CreateTime  time.Time   `gorm:"not null;autoCreateTime"`
	// ObjectNum   uint32      `gorm:"not null"`

	User User `gorm:"foreignKey:UserName;references:Name"`
}

func init() {
	if err := database.Client.AutoMigrate(&Bucket{}); err != nil {
		panic(err)
	}
}

func (aclTypeUtil) FromString(str string) ACLType {
	switch str {
	case "private":
		return PRIVATE
	case "public_read":
		return PUBLIC_READ
	case "public":
		return PUBLIC
	}
	return -1
}

func (aclTypeUtil) ToString(acl ACLType) string {
	switch acl {
	case PRIVATE:
		return "private"
	case PUBLIC_READ:
		return "public_read"
	case PUBLIC:
		return "public"
	}
	return "unknown"
}

func (storageTypeUtil) FromString(str string) StorageType {
	switch str {
	case "standard":
		return STANDARD
	case "low_freq":
		return LOW_FREQ
	}
	return -1
}

func (storageTypeUtil) ToString(storage_t StorageType) string {
	switch storage_t {
	case STANDARD:
		return "standard"
	case LOW_FREQ:
		return "low_freq"
	}
	return "unknown"
}
