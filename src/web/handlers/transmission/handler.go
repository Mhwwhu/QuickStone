package transmission

import (
	"io"
	"net/http"
	"strconv"

	"QuickStone/src/common"
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/webModels"
	meta "QuickStone/src/rpc/metadata"
	trans "QuickStone/src/rpc/transmission"
	grpcutil "QuickStone/src/utils/grpc"
	"QuickStone/src/web/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

var transClient trans.TransmissionServiceClient
var metaClient meta.MetadataServiceClient

func init() {
	conn := grpcutil.Connect(config.TransmissionServerName)
	transClient = trans.NewTransmissionServiceClient(conn)
	conn2 := grpcutil.Connect(config.MetadataServerName)
	metaClient = meta.NewMetadataServiceClient(conn2)
}

func UploadObjectHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)
	md, _ := metadata.FromOutgoingContext(ctx)
	userName := md.Get(constant.CtxUserNameKey)[0]
	userId, _ := strconv.Atoi(md.Get(constant.CtxUserIdKey)[0])

	var req webModels.UploadObjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.UploadObjectResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
			ObjectSize: 0,
		})
		return
	}

	if req.TargetUserName == "" {
		req.TargetUserName = userName
	}
	if req.Bucket == "" {
		req.Bucket = "Default"
	}
	if req.Key == "" {
		req.Key = uuid.New().String()
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		logrus.Fatalf("Cannot get file from request: %v", err)
		c.JSON(http.StatusOK, webModels.UploadObjectResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
			ObjectSize: 0,
		})
		return
	}

	/*
		TODO: 1. 向auth模块申请鉴权
	*/

	// 2. 向stroage-meta模块发送预存储请求
	{
		res, err := metaClient.RegisterUploadingObject(ctx, &meta.RegisterUploadingObjectRequest{
			TargetUserName: req.TargetUserName,
			Bucket:         req.Bucket,
			Key:            req.Key,
		})
		if err != nil || res.StatusCode != 0 {
			logrus.WithFields(logrus.Fields{
				"UserId": userId,
				"Bucket": req.Bucket,
				"Key":    req.Key,
			}).Warnf("Failed to register uploading: error=%v", err)
			c.JSON(http.StatusOK, webModels.UploadObjectResponse{
				StandardResponse: webModels.StandardResponse{
					StatusCode: res.StatusCode,
					StatusMsg:  res.StatusMsg,
				},
				ObjectSize: common.ObjectSizeT(fileHeader.Size),
			})
			return
		}
	}

	// 3. 向transmission模块流式发送存储请求
	{
		stream, err := transClient.UploadObject(ctx)
		if err != nil {
			logrus.Fatalf("Failed to get stream from Transmission server: %v", err)
		}

		stream.Send(&trans.UploadObjectRequestChunk{
			Payload: &trans.UploadObjectRequestChunk_Header{
				Header: &trans.UploadObjectRequestHeader{
					TargetUserName: req.TargetUserName,
					Bucket:         req.Bucket,
					Key:            req.Key,
					ObjectType:     req.ObjectType,
					ObjectSize:     uint64(fileHeader.Size),
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
		c.JSON(http.StatusOK, webModels.UploadObjectResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
			ObjectSize: common.ObjectSizeT(fileHeader.Size),
		})
		return
	}
}
