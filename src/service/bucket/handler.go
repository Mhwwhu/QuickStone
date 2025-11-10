package main

import (
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/rpc/bucket"
	"QuickStone/src/storage/database"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type BucketService struct {
	bucket.BucketServiceServer
}

func (s BucketService) Init() {

}

func (s BucketService) CreateBucket(ctx context.Context, req *bucket.CreateBucketRequest) (resp *bucket.CreateBucketResponse, err error) {
	md, _ := metadata.FromIncomingContext(ctx)

	resp = &bucket.CreateBucketResponse{}
	var bucketModel dbModels.Bucket
	result := database.Client.WithContext(ctx).Limit(1).Where("name = ?", req.Bucket).Find(&bucketModel)
	if result.RowsAffected != 0 {
		logrus.Infof("Bucket %s already exists.", req.Bucket)
		resp = &bucket.CreateBucketResponse{
			StatusCode: constant.BucketExistsErrorCode,
			StatusMsg:  constant.BucketExistsError,
		}
		return
	}

	bucketModel.Name = req.Bucket
	bucketModel.Area = req.Area
	bucketModel.StorageType = dbModels.StorageType(req.StorageType)
	bucketModel.ACLType = dbModels.ACLType(req.AclType)
	bucketModel.UserName = md.Get(constant.CtxUserNameKey)[0]
	result = database.Client.WithContext(ctx).Create(&bucketModel)
	if result.Error != nil {
		logrus.Errorf("Database operation failed: %v", result.Error)
		resp = &bucket.CreateBucketResponse{
			StatusCode: constant.DatabaseErrorCode,
		}
		return
	}

	resp = &bucket.CreateBucketResponse{
		StatusCode:      0,
		StatusMsg:       "",
		CreateTimestamp: bucketModel.CreateTime.Format(config.TimeFormat),
	}
	logrus.WithFields(logrus.Fields{
		"UserName": bucketModel.UserName,
		"Bucket":   bucketModel.Name,
	}).Infof("The bucket has been created.")
	return
}

func (s BucketService) DeleteBucket(ctx context.Context, req *bucket.DeleteBucketRequest) (resp *bucket.DeleteBucketResponse, err error) {
	return nil, nil
}

func (s BucketService) ShowBucket(ctx context.Context, req *bucket.ShowBucketRequest) (resp *bucket.ShowBucketResponse, err error) {
	resp = &bucket.ShowBucketResponse{}
	var bucketModel dbModels.Bucket
	result := database.Client.WithContext(ctx).Where("name = ?", req.Bucket).Find(&bucketModel)
	if result.RowsAffected == 0 {
		logrus.Infof("No matched buckets.")
		resp.StatusCode = constant.BucketNotExistsErrorCode
		resp.StatusMsg = constant.BucketNotExistsError
		return
	}

	resp.Area = bucketModel.Area
	resp.AclType = bucket.BucketACLType(bucketModel.ACLType)
	resp.StorageType = bucket.StorageType(bucketModel.StorageType)
	resp.CreateTimestamp = bucketModel.CreateTime.Format(config.TimeFormat)

	var count int64
	database.Client.WithContext(ctx).Where("user_name = ? and bucket_name = ?", req.UserName, req.Bucket).Count(&count)
	resp.ObjectNum = uint32(count)

	return
}
