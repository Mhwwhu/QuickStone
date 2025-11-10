package main

import (
	"context"
	"fmt"

	"QuickStone/src/common"
	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/rpc/metadata"
	"QuickStone/src/storage/database"
	"QuickStone/src/utils/cache"
	"QuickStone/src/utils/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type MetadataService struct {
	metadata.MetadataServiceServer
}

var conn *amqp.Connection
var channel *amqp.Channel

func (MetadataService) Init() {
	conn = rabbitmq.ConnectMQ()
	var err error
	channel, err = conn.Channel()
	common.ExitOnErr(err)
}

func CloseMQConn() {
	if err := conn.Close(); err != nil {
		panic(err)
	}

	if err := channel.Close(); err != nil {
		panic(err)
	}
}

func (MetadataService) RegisterUploadingObject(ctx context.Context, req *metadata.RegisterUploadingObjectRequest) (
	resp *metadata.RegisterUploadingObjectResponse, err error) {
	key := fmt.Sprintf("%s:register_upload:%s:%s:%s", constant.MetadataVarPrefix, req.TargetUserName, req.Bucket, req.Key)
	registered := cache.Get[bool](ctx, key)

	// 先检查缓存里面有没有，如果没有则读数据库
	if registered != nil && *registered {
		return &metadata.RegisterUploadingObjectResponse{
			StatusCode: constant.ObjectUploadConflictErrorCode,
			StatusMsg:  constant.ObjectUploadConflictError,
		}, nil
	}

	var count int64
	result := database.Client.WithContext(ctx).Model(&dbModels.Object{}).
		Where("user_name = ? and bucket_name = ? and key = ?", req.TargetUserName, req.Bucket, req.Key).
		Count(&count)
	if result.Error != nil {
		logrus.Errorf("Error on database: %v", result.Error)
		return &metadata.RegisterUploadingObjectResponse{
			StatusCode: constant.DatabaseErrorCode,
			StatusMsg:  "",
		}, nil
	}
	if count != 0 {
		return &metadata.RegisterUploadingObjectResponse{
			StatusCode: constant.ObjectUploadConflictErrorCode,
			StatusMsg:  constant.ObjectUploadConflictError,
		}, nil
	}

	cache.Set(ctx, key, true)
	return &metadata.RegisterUploadingObjectResponse{StatusCode: 0, StatusMsg: ""}, nil
}
