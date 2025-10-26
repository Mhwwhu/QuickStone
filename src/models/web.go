package models

import "QuickStone/src/common"

// 对接客户端和网关的数据模型

type LoginRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	StatusCode common.StatusCodeT `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	UserId     common.UserIdT     `json:"user_id"`
	Token      string             `json:"token"`
}

type RegisterRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterResponse struct {
	StatusCode common.StatusCodeT `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	UserId     common.UserIdT     `json:"user_id"`
	Token      string             `json:"token"`
}

type UploadObjectRequest struct {
	TargetUserId common.UserIdT `form:"target_user_id"` // 留空表示当前用户
	Bucket       string         `form:"bucket" binding:"required"`
	Key          string         `form:"key" binding:"required"`
	ObjectType   string         `form:"obj_type"`
}

type UploadObjectResponse struct {
	StatusCode common.StatusCodeT `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	ObjectSize common.ObjectSizeT `json:"object_size"`
}

type DownloadObjectRequest struct {
}
