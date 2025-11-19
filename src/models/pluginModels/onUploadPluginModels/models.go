package onuploadpluginmodels

import (
	"QuickStone/src/common"
	"time"
)

type ObjectModel struct {
	UserName   string             `json:"user_name"`
	BucketName string             `json:"bucket_name"`
	Key        string             `json:"key"`
	ObjectType string             `json:"object_type"`
	Size       common.ObjectSizeT `json:"size"`
	CreateTime time.Time          `json:"create_time"`
}

type OnUploadContext struct {
	UserName string `json:"user_name"`
	UserId   string `json:"user_id"`

	ObjectMeta ObjectModel `json:"object_meta"`
	Object     []byte      `json:"object"`
}

type OnUploadAction struct {
	ActionType string                 `json:"action_type"`
	Payload    map[string]interface{} `json:"payload"` // 操作的具体负载，由于操作不同，负载也可能不一样
}
