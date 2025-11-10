package common

import "strconv"

type StatusCodeT = uint32
type UserIdT = uint32
type ObjectSizeT = uint64

type CtxKeyT string

type StoragePath struct {
	UserName string
	Bucket   string
	Key      string
}

func AtoUserIdT(str string) UserIdT {
	userIdInt, _ := strconv.Atoi(str)
	return UserIdT(userIdInt)
}

func ExitOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
