package transmission

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	"github.com/mhwwhu/QuickStone/src/constant"
	"github.com/mhwwhu/QuickStone/src/models"
	"github.com/mhwwhu/QuickStone/src/rpc/metadata"
	trans "github.com/mhwwhu/QuickStone/src/rpc/transmission"
	grpcutil "github.com/mhwwhu/QuickStone/src/utils/grpc"
	"github.com/sirupsen/logrus"
)

var TransClient trans.TransmissionServiceClient
var MetaClient metadata.MetadataServiceClient

func init() {
	conn := grpcutil.Connect(config.TransmissionServerName)
	TransClient = trans.NewTransmissionServiceClient(conn)
}

func UploadObjectHandle(c *gin.Context) {
	userId := c.GetUint64(constant.CtxUserIdKey)
	ctx := context.WithValue(c.Request.Context(), constant.CtxUserIdKey, userId)
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
		2. 向stroage-meta模块发送预存储请求
		3. 向transmission模块发送存储请求
	*/

	// 2. 向stroage-meta模块发送预存储请求
	{
		res, err := MetaClient.RegisterUploadingObject(ctx, &metadata.RegisterUploadingObjectRequest{
			TargetUserId: req.TargetUserId,
			Bucket:       req.Bucket,
			Key:          req.Key,
			UserId:       userId,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"UserId": userId,
				"Bucket": req.Bucket,
				"Key":    req.Key,
			}).Warnf("Error when trying to register uploading: %v", err)
			c.JSON(http.StatusOK, models.UploadObjectResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			})
			return
		}
	}

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

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"UserId": userId,
				"Bucket": req.Bucket,
				"Key":    req.Key,
			}).Warnf("Error when trying to upload: %v", err)
			c.JSON(http.StatusOK, models.UploadObjectResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
