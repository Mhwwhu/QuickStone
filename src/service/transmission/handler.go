package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/msgModels"
	trans "QuickStone/src/rpc/transmission"
	"QuickStone/src/storage"
	"QuickStone/src/utils/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sirupsen/logrus"
)

var conn *amqp.Connection
var channel *amqp.Channel

type TransmissionService struct {
	trans.TransmissionServiceServer
}

func (TransmissionService) Init() {
	conn = rabbitmq.ConnectMQ()
	channel, _ = conn.Channel()
	channel.ExchangeDeclare(
		constant.ObjectStorageExchange, // exchange 名
		"topic",                        // 类型
		true,                           // durable：持久化
		false,                          // auto-deleted：不自动删除
		false,                          // internal：客户端可用
		false,                          // no-wait
		nil,                            // arguments
	)

	_, err := channel.QueueDeclare(
		constant.ObjectMetaQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	common.ExitOnErr(err)

	err = channel.QueueBind(
		constant.ObjectMetaQueue,
		fmt.Sprintf("%s.*", constant.ObjectEventPrefix),
		constant.ObjectStorageExchange,
		false,
		nil,
	)
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

func (s TransmissionService) UploadObject(stream trans.TransmissionService_UploadObjectServer) error {
	req, err := stream.Recv()
	if err != nil {
		stream.SendAndClose(&trans.UploadObjectResponse{
			StatusCode: constant.GrpcCommunicationErrorCode,
		})
		return err
	}

	head := req.GetHeader()
	storagePath := common.StoragePath{
		UserName: head.TargetUserName,
		Bucket:   head.Bucket,
		Key:      head.Key,
	}

	var buf bytes.Buffer
	frameNum := uint32((head.ObjectSize-1)/config.GrpcStreamUploadSliceSize + 1)
	frameNo := uint32(0)
	for frameNo < frameNum {
		req, err := stream.Recv()
		body := req.GetData()
		if err != nil || frameNo != body.SeriesNo {
			logrus.Errorf("Grpc communication error: err = %v, frameNo = %d, seriesNo = %d", err, frameNo, body.SeriesNo)
			stream.SendAndClose(&trans.UploadObjectResponse{
				StatusCode: constant.GrpcCommunicationErrorCode,
			})
			return err
		}
		frameNo++

		data := body.Data
		if _, err := buf.Write(data); err != nil {
			stream.SendAndClose(&trans.UploadObjectResponse{
				StatusCode: constant.GolangInternalErrorCode,
			})
			return err
		}
	}

	reader := bytes.NewReader(buf.Bytes())

	// TODO: 目前是同步等待存储模块，后续可以改成MQ异步发布任务
	err = storage.StorageClient.UploadObject(stream.Context(), storagePath, reader)
	if err != nil {
		stream.SendAndClose(&trans.UploadObjectResponse{
			StatusCode: constant.StorageUploadErrorCode,
		})
		return err
	}

	stream.SendAndClose(&trans.UploadObjectResponse{
		StatusCode: 0,
	})

	// 发布对象存储消息
	event := msgModels.Object{
		EventType:  "stored",
		UserName:   head.TargetUserName,
		Bucket:     head.Bucket,
		Key:        head.Key,
		ObjType:    head.ObjectType,
		Size:       head.ObjectSize,
		OccurredAt: time.Now(),
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel.PublishWithContext(
		stream.Context(),
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

	return nil
}

func (s TransmissionService) DownloadObject(ctx context.Context, req *trans.DownloadObjectRequest) (*trans.DownloadObjectResponse, error) {
	return nil, nil
}
