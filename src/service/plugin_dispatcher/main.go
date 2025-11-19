package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"QuickStone/src/common"
	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/models/msgModels"
	"QuickStone/src/storage"
	"QuickStone/src/storage/database"
	"QuickStone/src/utils/cache"

	"github.com/oklog/run"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {
	// go updateObjMeta(channel)

	g := &run.Group{}
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
		constant.ObjectOnUploadProcessQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	for msg := range msgs {
		logrus.Info("000\n")
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
			logrus.Infof("222\n")
			result := database.Client.Unscoped().
				Where("user_name = ? and bucket_name = ? and key = ?", evt.UserName, evt.Bucket, evt.Key).
				Delete(&dbModels.Object{})
			common.ExitOnErr(result.Error)
		}
		msg.Ack(false)
	}
}

func deleteObjectHard(channel *amqp.Channel) {

}
