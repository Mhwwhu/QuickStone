package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"QuickStone/src/common"
	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/models/msgModels"
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
	// registered := cache.Get[bool](ctx, key)

	// 先检查缓存里面有没有，如果没有则读数据库
	// if registered != nil && *registered {
	// 	return &metadata.RegisterUploadingObjectResponse{
	// 		StatusCode: constant.ObjectUploadConflictErrorCode,
	// 		StatusMsg:  constant.ObjectUploadConflictError,
	// 	}, nil
	// }

	var count int64
	result := database.Client.WithContext(ctx).Model(&dbModels.Object{}).
		Where("user_name = ? and bucket_name = ? and key = ? and is_deleted = ?", req.TargetUserName, req.Bucket, req.Key, false).
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

func (MetadataService) DeleteObject(ctx context.Context, req *metadata.DeleteObjectRequest) (resp *metadata.DeleteObjectResponse, err error) {
	// TODO: 使用缓存
	var objectModel dbModels.Object
	result := database.Client.WithContext(ctx).Model(&dbModels.Object{}).
		Where("user_name = ? and bucket_name = ? and key = ?", req.TargetUserName, req.Bucket, req.Key).
		Find(&objectModel)
	if result.Error != nil {
		logrus.Errorf("Error on database: %v", result.Error)
		return &metadata.DeleteObjectResponse{
			StatusCode: constant.DatabaseErrorCode,
			StatusMsg:  "",
		}, nil
	}
	if result.RowsAffected == 0 || objectModel.IsDeleted {
		return &metadata.DeleteObjectResponse{
			StatusCode: constant.ObjectDeletedErrorCode,
			StatusMsg:  constant.ObjectDeletedError,
		}, nil
	}

	// 软删除，等物理存储删除之后再硬删除
	result = database.Client.WithContext(ctx).Model(&dbModels.Object{}).
		Where("user_name = ? and bucket_name = ? and key = ?", req.TargetUserName, req.Bucket, req.Key).
		Update("is_deleted", true)
	if result.Error != nil {
		logrus.Errorf("Error on database: %v", result.Error)
		return &metadata.DeleteObjectResponse{
			StatusCode: constant.DatabaseErrorCode,
			StatusMsg:  "",
		}, nil
	}
	logrus.Infof("888\n")
	event := msgModels.Object{
		EventType:  "soft-deleted",
		UserName:   req.TargetUserName,
		Bucket:     req.Bucket,
		Key:        req.Key,
		OccurredAt: time.Now(),
	}
	body, err := json.Marshal(event)
	if err != nil {
		return &metadata.DeleteObjectResponse{
			StatusCode: constant.InternalErrorCode,
			StatusMsg:  "",
		}, err
	}

	channel.PublishWithContext(
		ctx,
		constant.ObjectStorageExchange,
		constant.ObjectStoredEvent,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	return &metadata.DeleteObjectResponse{}, nil
}
