package transmission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	"github.com/mhwwhu/QuickStone/src/constant"
	"github.com/mhwwhu/QuickStone/src/models"
	trans "github.com/mhwwhu/QuickStone/src/rpc/transmission"
	grpcutil "github.com/mhwwhu/QuickStone/src/utils/grpc"
)

var Client trans.TransmissionServiceClient

func init() {
	conn := grpcutil.Connect(config.TransmissionServerName)
	Client = trans.NewTransmissionServiceClient(conn)
}

func UploadObjectHandle(c *gin.Context) {
	var req models.UploadObjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.UploadObjectResponse{
			StatusCode: constant.GateWayParamsErrorCode,
			StatusMsg:  constant.GateWayParamsError,
			ObjectSize: 0,
		})
		return
	}

	// TODO: 向auth模块申请鉴权

	// res, err := Client.UploadObject(c.Request.Context(), &trans.UploadObjectRequest{
	// 	Token:      req.Token,
	// 	UserId:     req.UserId,
	// 	Bucket:     req.Bucket,
	// 	Key:        req.Key,
	// 	ObjectType: req.ObjectType,
	// })

	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"Username": req.,
	// 	}).Warnf("Error when trying to upload: %v", err)
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }

	// logrus.WithFields(logrus.Fields{
	// 	"Username": req.UserName,
	// 	"Token":    res.Token,
	// 	"UserId":   res.Uid,
	// }).Infof("User %s has log in", req.UserName)

	// c.JSON(http.StatusOK, res)
}
