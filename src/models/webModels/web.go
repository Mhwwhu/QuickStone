package webModels

import "QuickStone/src/common"

// 对接客户端和网关的数据模型

type StandardResponse struct {
	StatusCode common.StatusCodeT `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginResponse struct {
	StandardResponse
	UserId common.UserIdT `json:"user_id"`
	Token  string         `json:"token"`
}

type RegisterRequest struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterResponse struct {
	StandardResponse
	UserId common.UserIdT `json:"user_id"`
	Token  string         `json:"token"`
}

type UploadObjectRequest struct {
	TargetUserName string `json:"target_user_name" form:"target_user_name"` // 留空表示当前用户
	Bucket         string `json:"bucket_name" form:"bucket_name"`
	Key            string `json:"key" form:"key"`
	ObjectType     string `json:"object_type" form:"object_type"`
}

type UploadObjectResponse struct {
	StandardResponse
	ObjectSize common.ObjectSizeT `json:"object_size"`
}

type DownloadObjectRequest struct {
}

type CreateBucketRequest struct {
	Bucket      string `json:"bucket_name" form:"bucket_name" binding:"required"`
	Area        string `json:"area" form:"area"`
	StorageType string `json:"storage_type" form:"storage_type"`
	ACLType     string `json:"acl_type" form:"acl_type"`
}

type CreateBucketResponse struct {
	StandardResponse
	CreateTime string `json:"create_time"`
}

type ShowBucketRequest struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Bucket   string `json:"bucket_name" form:"bucket_name" binding:"required"`
}

type ShowBucketResponse struct {
	StandardResponse
	Area        string `json:"area"`
	StorageType string `json:"storage_type"`
	ACLType     string `json:"acl_type"`
	CreateTime  string `json:"create_time"`
	ObjectNum   uint32 `json:"object_num"`
	Status      string `json:"status"`
}

type ShowUserBucketsRequest struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}

type BucketMeta struct {
	UserName    string `json:"user_name"`
	BucketName  string `json:"bucket_name"`
	Area        string `json:"area"`
	StorageType string `json:"storage_type"`
	ACLType     string `json:"acl_type"`
	CreateTime  string `json:"create_time"`
	Status      string `json:"status"`
}

type ShowUserBucketsResponse struct {
	StandardResponse
	Buckets []BucketMeta `json:"buckets"`
}

type ShowObjectsRequest struct {
	UserName   string `json:"user_name" form:"user_name"`
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
}

type ObjectMeta struct {
	Key        string             `json:"key"`
	Size       common.ObjectSizeT `json:"size"`
	CreateTime string             `json:"create_time"`
}

type ShowObjectsResponse struct {
	StandardResponse
	Objects []ObjectMeta `json:"objects"`
}
