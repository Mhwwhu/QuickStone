package models

type LoginRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint32 `json:"user_id"`
	Token      string `json:"token"`
}

type RegisterRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint32 `json:"user_id"`
	Token      string `json:"token"`
}

type UploadObjectRequest struct {
	Token      string `form:"token" binding:"required"`
	UserId     string `form:"user_id"` // 留空表示当前用户
	Bucket     string `form:"bucket" binding:"required"`
	Key        string `form:"key" binding:"required"`
	ObjectType string `form:"obj_type"`
}

type UploadObjectResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	ObjectSize int    `json:"object_size"`
}

type DownloadObjectRequest struct {
}
