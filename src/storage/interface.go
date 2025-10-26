package storage

import (
	"context"
	"io"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/storage/fs"
)

var StorageClient IStorage

func init() {
	switch config.EnvCfg.StorageType {
	case "fs":
		StorageClient = fs.FSStorage{}
	}
}

type IStorage interface {
	UploadObject(ctx context.Context, path common.StoragePath, reader io.Reader) error
}
