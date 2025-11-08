package models

import "QuickStone/src/common"

// 对接客户端和网关的数据模型

type StandardResponse struct {
	StatusCode common.StatusCodeT `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
}

type LoginRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	StandardResponse
	UserId common.UserIdT `json:"user_id"`
	Token  string         `json:"token"`
}

type RegisterRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterResponse struct {
	StandardResponse
	UserId common.UserIdT `json:"user_id"`
	Token  string         `json:"token"`
}

type UploadObjectRequest struct {
	TargetUserName string `form:"target_user_name"` // 留空表示当前用户
	Bucket         string `form:"bucket"`
	Key            string `form:"key"`
	ObjectType     string `form:"obj_type"`
}

type UploadObjectResponse struct {
	StandardResponse
	ObjectSize common.ObjectSizeT `json:"object_size"`
}

type DownloadObjectRequest struct {
}

type CreateBucketRequest struct {
	Bucket      string `form:"bucket" binding:"required"`
	Area        string `form:"area"`
	StorageType string `form:"storage_type"`
	ACLType     string `form:"acl_type"`
}

type CreateBucketResponse struct {
	StandardResponse
	CreateTime string `json:"create_time"`
}

type ShowBucketRequest struct {
	UserName string `form:"user_name" binding:"required"`
	Bucket   string `form:"bucket" binding:"required"`
}

type ShowBucketResponse struct {
	StandardResponse
	Area        string `json:"area"`
	StorageType string `json:"storage_type"`
	ACLType     string `json:"acl_type"`
	CreateTime  string `json:"create_time"`
	Status      string `json:"status"`
}
