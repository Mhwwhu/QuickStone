package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"syscall"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/models/msgModels"
	"QuickStone/src/rpc/metadata"
	"QuickStone/src/storage"
	"QuickStone/src/storage/database"
	"QuickStone/src/utils/cache"
	"QuickStone/src/utils/consul"

	"github.com/oklog/run"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var ServerId string

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.EnvCfg.PodIpAddr, config.MetadataServerPort))

	if err != nil {
		logrus.Panicf("Rpc %s listen error: %v", config.MetadataServerName, err)
	}

	s := grpc.NewServer()

	if ServerId, err = consul.RegisterConsul(config.MetadataServerName, config.MetadataServerPort); err != nil {
		logrus.Panicf("Rpc %s registering consul error: %v", config.MetadataServerName, err)
	}

	var srv MetadataService
	metadata.RegisterMetadataServiceServer(s, srv)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	srv.Init()

	// defer CloseMQConn()

	go updateObjMeta(channel)

	g := &run.Group{}
	g.Add(func() error {
		return s.Serve(lis)
	}, func(err error) {
		s.GracefulStop()
		s.Stop()
		logrus.Errorf("Rpc %s listening error: %v", config.TransmissionServerName, err)
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when runing http server")
		os.Exit(1)
	}
}

func updateObjMeta(channel *amqp.Channel) {
	msgs, _ := channel.Consume(
		constant.ObjectMetaQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	for msg := range msgs {
		evt := msgModels.Object{}
		json.Unmarshal(msg.Body, &evt)
		switch evt.EventType {
		case "stored":
			obj := dbModels.Object{
				UserName:   evt.UserName,
				BucketName: evt.Bucket,
				Key:        evt.Key,
				ObjectType: evt.ObjType,
				Size:       evt.Size,
				IsDeleted:  false,
			}
			result := database.Client.Create(&obj)
			if result.Error != nil {
				cache.Set(
					context.Background(),
					fmt.Sprintf("%s:register_upload:%s:%s:%s", constant.MetadataVarPrefix, evt.UserName, evt.Bucket, evt.Key),
					false,
				)
			}
		case "soft-deleted":
			// 调用存储模块进行物理删除
			logrus.Infof("111\n")
			err := storage.StorageClient.DeleteObject(
				context.Background(),
				common.StoragePath{
					UserName: evt.UserName,
					Bucket:   evt.Bucket,
					Key:      evt.Key,
				},
			)
			common.ExitOnErr(err)
			evt.EventType = "physical-deleted"
			body, err := json.Marshal(evt)
			common.ExitOnErr(err)
			channel.PublishWithContext(
				context.Background(),
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
		case "physical-deleted":
			result := database.Client.Unscoped().
				Where("user_name = ? and bucket_name = ? and key = ?", evt.UserName, evt.Bucket, evt.Key).
				Delete(&dbModels.Object{})
			common.ExitOnErr(result.Error)
		}
		msg.Ack(false)
	}
}
