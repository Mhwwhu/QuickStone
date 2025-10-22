package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	"github.com/mhwwhu/QuickStone/src/constant"
	"github.com/mhwwhu/QuickStone/src/models"
	"github.com/mhwwhu/QuickStone/src/rpc/auth"
	grpcutil "github.com/mhwwhu/QuickStone/src/utils/grpc"
	"github.com/sirupsen/logrus"
)

var Client auth.AuthServiceClient

func init() {
	conn := grpcutil.Connect(config.AuthServerName)
	Client = auth.NewAuthServiceClient(conn)
}

func LoginHandle(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.LoginResponse{
			StatusCode: constant.GateWayParamsErrorCode,
			StatusMsg:  constant.GateWayParamsError,
			UserId:     0,
			Token:      "",
		})
		return
	}

	res, err := Client.Login(c.Request.Context(), &auth.LoginRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to login: %v", err)
		c.JSON(http.StatusOK, res)
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.Uid,
	}).Infof("User %s has log in", req.UserName)

	c.JSON(http.StatusOK, res)
}

func RegisterHandle(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, auth.LoginResponse{
			StatusCode: constant.GateWayParamsErrorCode,
			StatusMsg:  constant.GateWayParamsError,
			Uid:        0,
			Token:      "",
		})
		return
	}

	res, err := Client.Register(c.Request.Context(), &auth.RegisterRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to register")
		c.JSON(http.StatusOK, res)
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.Uid,
	}).Infof("User %s has registered.", req.UserName)

	c.JSON(http.StatusOK, res)
}
