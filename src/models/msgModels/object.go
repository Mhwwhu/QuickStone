package msgModels

import "time"

type Object struct {
	EventType  string    `json:"eventType"` // "stored" / "updated" / "deleted"
	UserName   string    `json:"userName"`
	Bucket     string    `json:"bucket"`
	Key        string    `json:"key"`
	Size       uint64    `json:"size"`
	ObjType    string    `json:"objType"`
	OccurredAt time.Time `json:"occurredAt"`
}
