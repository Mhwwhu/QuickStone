package metadata

import (
	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/webModels"
	meta "QuickStone/src/rpc/metadata"
	grpcutil "QuickStone/src/utils/grpc"
	"QuickStone/src/web/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var metaClient meta.MetadataServiceClient

func init() {
	conn := grpcutil.Connect(config.MetadataServerName)
	metaClient = meta.NewMetadataServiceClient(conn)
}

func DeleteObjectHandle(c *gin.Context) {
	ctx := utils.CreateCtxFromGin(c)
	md, _ := metadata.FromOutgoingContext(ctx)
	userName := md.Get(constant.CtxUserNameKey)[0]

	var req webModels.DeleteObjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.DeleteObjectResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
		})
		return
	}

	if req.TargetUserName == "" {
		req.TargetUserName = userName
	}
	if req.Bucket == "" {
		req.Bucket = "Default"
	}

	resp, err := metaClient.DeleteObject(ctx, &meta.DeleteObjectRequest{
		TargetUserName: req.TargetUserName,
		Bucket:         req.Bucket,
		Key:            req.Key,
	})

	if resp.StatusCode != 0 || err != nil {
		c.JSON(http.StatusOK, webModels.DeleteObjectResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: resp.StatusCode,
				StatusMsg:  resp.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, webModels.DeleteObjectResponse{})
}
