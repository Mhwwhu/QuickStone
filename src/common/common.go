package common

type StatusCodeT = uint32
type UserIdT = uint32
type ObjectSizeT = uint64

type CtxKeyT string

type StoragePath struct {
	UserId UserIdT
	Bucket string
	Key    string
}
