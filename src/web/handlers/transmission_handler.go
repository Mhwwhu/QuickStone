package handlers

import (
	"context"
	"io"
	"net/http"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models"
	"QuickStone/src/rpc/metadata"
	trans "QuickStone/src/rpc/transmission"
	grpcutil "QuickStone/src/utils/grpc"
	"QuickStone/src/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var TransClient trans.TransmissionServiceClient
var MetaClient metadata.MetadataServiceClient

func init() {
	conn := grpcutil.Connect(config.TransmissionServerName)
	TransClient = trans.NewTransmissionServiceClient(conn)
}

func UploadObjectHandle(c *gin.Context) {
	claim, _ := c.Get(constant.CtxClaimKey)
	userId := claim.(jwt.Claims).UserID
	userName := claim.(jwt.Claims).Username
	ctx := context.WithValue(c.Request.Context(), constant.CtxUserIdKey, userId)
	ctx = context.WithValue(ctx, constant.CtxUserNameKey, userName)

	var req models.UploadObjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.UploadObjectResponse{
			StatusCode: constant.GateWayParamsErrorCode,
			StatusMsg:  constant.GateWayParamsError,
			ObjectSize: 0,
		})
		return
	}

	if req.TargetUserId == 0 {
		req.TargetUserId = userId
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		logrus.Fatalf("Cannot get file from request: %v", err)
		c.JSON(http.StatusOK, models.UploadObjectResponse{
			StatusCode: constant.GateWayParamsErrorCode,
			StatusMsg:  constant.GateWayParamsError,
			ObjectSize: 0,
		})
		return
	}

	/*
		TODO: 1. 向auth模块申请鉴权
	*/

	// 2. 向stroage-meta模块发送预存储请求
	// {
	// 	res, err := MetaClient.RegisterUploadingObject(ctx, &metadata.RegisterUploadingObjectRequest{
	// 		TargetUserId: req.TargetUserId,
	// 		Bucket:       req.Bucket,
	// 		Key:          req.Key,
	// 		UserId:       userId,
	// 	})
	// 	if err != nil {
	// 		logrus.WithFields(logrus.Fields{
	// 			"UserId": userId,
	// 			"Bucket": req.Bucket,
	// 			"Key":    req.Key,
	// 		}).Warnf("Error when trying to register uploading: %v", err)
	// 		c.JSON(http.StatusOK, models.UploadObjectResponse{
	// 			StatusCode: res.StatusCode,
	// 			StatusMsg:  res.StatusMsg,
	// 		})
	// 		return
	// 	}
	// }

	// 3. 向transmission模块流式发送存储请求
	{
		stream, err := TransClient.UploadObject(ctx)
		if err != nil {
			logrus.Fatalf("Failed to get stream from Transmission server: %v", err)
		}

		stream.Send(&trans.UploadObjectRequestChunk{
			Payload: &trans.UploadObjectRequestChunk_Header{
				Header: &trans.UploadObjectRequestHeader{
					TargetUserId: req.TargetUserId,
					Bucket:       req.Bucket,
					Key:          req.Key,
					ObjectType:   req.ObjectType,
					ObjectSize:   uint64(fileHeader.Size),
				},
			},
		})

		var seriesNo uint32 = 0
		for {
			buf := make([]byte, config.GrpcStreamUploadSliceSize)
			size, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			stream.Send(&trans.UploadObjectRequestChunk{
				Payload: &trans.UploadObjectRequestChunk_Data{
					Data: &trans.UploadObjectRequestBody{
						SeriesNo: seriesNo,
						Data:     buf[:size],
					},
				},
			})
			seriesNo++
		}

		res, err := stream.CloseAndRecv()
		logrus.Infof("received resp: %d", res.StatusCode)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"UserId": userId,
				"Bucket": req.Bucket,
				"Key":    req.Key,
			}).Warnf("Error when trying to upload: %v", err)
		}
		c.JSON(http.StatusOK, models.UploadObjectResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
			ObjectSize: common.ObjectSizeT(fileHeader.Size),
		})
		return
	}
}
