package bucket

import (
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models"
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

	var req models.CreateBucketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.CreateBucketResponse{
			StandardResponse: models.StandardResponse{
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
		c.JSON(http.StatusOK, models.CreateBucketResponse{
			StandardResponse: models.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.CreateBucketResponse{
		StandardResponse: models.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		CreateTime: res.CreateTimestamp,
	})
}

func ShowBucketHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)

	var req models.ShowBucketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.ShowBucketResponse{
			StandardResponse: models.StandardResponse{
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
		c.JSON(http.StatusOK, models.ShowBucketResponse{
			StandardResponse: models.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.ShowBucketResponse{
		StandardResponse: models.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		CreateTime:  res.CreateTimestamp,
		Area:        res.Area,
		ACLType:     bucket.BucketACLTypeUtil.ToString(res.AclType),
		StorageType: bucket.StorageTypeUtil.ToString(res.StorageType),
		Status:      "OK",
	})
}

func OverviewHandle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"bucket_list": "",
		"bucket_num":  9,
	})
}
