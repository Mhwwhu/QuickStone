package fs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"QuickStone/src/common"
	"QuickStone/src/config"
)

type FSStorage struct{}

func (FSStorage) UploadObject(ctx context.Context, path common.StoragePath, reader io.Reader) error {
	localPath := fmt.Sprintf("%s/%s/%s/%s", config.EnvCfg.FsRootPath, path.UserName, path.Bucket, path.Key)
	// 创建目标文件的所有目录
	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories for path %s: %w", localPath, err)
	}

	// 创建文件
	file, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file at path %s: %w", localPath, err)
	}
	defer file.Close()

	// 将数据从 reader 写入文件
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write data to file %s: %w", localPath, err)
	}

	return nil
}
