package bucket

import (
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/webModels"
	"QuickStone/src/rpc/bucket"
	grpcutil "QuickStone/src/utils/grpc"
	"QuickStone/src/web/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var bucketClient bucket.BucketServiceClient

func init() {
	conn := grpcutil.Connect(config.BucketServerName)
	bucketClient = bucket.NewBucketServiceClient(conn)
}

func CreateBucketHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)

	var req webModels.CreateBucketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.CreateBucketResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
		})
		return
	}

	res, err := bucketClient.CreateBucket(ctx, &bucket.CreateBucketRequest{
		Bucket:      req.Bucket,
		Area:        req.Area,
		AclType:     bucket.BucketACLTypeUtil.FromString(req.ACLType),
		StorageType: bucket.StorageTypeUtil.FromString(req.StorageType),
	})

	if err != nil {
		c.JSON(http.StatusOK, webModels.CreateBucketResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, webModels.CreateBucketResponse{
		StandardResponse: webModels.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		CreateTime: res.CreateTimestamp,
	})
}

func ShowBucketHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)

	var req webModels.ShowBucketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.ShowBucketResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
		})
		return
	}

	res, err := bucketClient.ShowBucket(ctx, &bucket.ShowBucketRequest{
		Bucket:   req.Bucket,
		UserName: req.UserName,
	})

	if err != nil {
		c.JSON(http.StatusOK, webModels.ShowBucketResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, webModels.ShowBucketResponse{
		StandardResponse: webModels.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		CreateTime:  res.CreateTimestamp,
		Area:        res.Area,
		ACLType:     bucket.BucketACLTypeUtil.ToString(res.AclType),
		StorageType: bucket.StorageTypeUtil.ToString(res.StorageType),
		ObjectNum:   res.ObjectNum,
		Status:      "OK",
	})
}

func ShowUserBucketsHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)

	var req webModels.ShowUserBucketsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.ShowUserBucketsResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
		})
		return
	}

	res, err := bucketClient.ShowUserBuckets(ctx, &bucket.ShowUserBucketsRequest{
		UserName: req.UserName,
	})

	if err != nil {
		c.JSON(http.StatusOK, webModels.ShowUserBucketsResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	resp := webModels.ShowUserBucketsResponse{
		StandardResponse: webModels.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
	}
	for _, b := range res.Buckets {
		resp.Buckets = append(resp.Buckets, webModels.BucketMeta{
			UserName:    req.UserName,
			BucketName:  b.BucketName,
			Area:        b.Area,
			StorageType: bucket.StorageType_name[int32(b.StorageType)],
			ACLType:     bucket.BucketACLType_name[int32(b.AclType)],
			CreateTime:  b.CreateTimestamp,
		})
	}
	c.JSON(http.StatusOK, resp)
}
