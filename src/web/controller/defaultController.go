package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/models"
	"github.com/sirupsen/logrus"
)

// var Client auth.AuthServiceClient
// func init() {
// 	conn := grpcutil.Connect(config.AuthServerName)
// 	Client = auth.NewAuthServiceClient(conn)
// }

var users = map[string]string{
	"123@qq.com": "123",
	"bob":        "123",
}

func LoginHandle(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 400,
			"status_msg":  "参数错误",
			"user_id":     0,
		})
		return
	}

	pwd, ok := users[req.UserName]
	if !ok || pwd != req.Password {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 401,
			"status_msg":  "用户名或密码错误",
			"user_id":     0,
		})
		return
	}

	// 登录成功，保存 session

	session := sessions.Default(c)
	session.Set("userid", req.UserName)
	session.Set("username", req.UserName)
	if err := session.Save(); err != nil {
		logrus.Errorf("Failed to save session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "保存 session 失败",
		})
		return
	}

	logrus.Infof("User %s has logged in", req.UserName)

	// 返回登录成功信息
	c.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"status_msg":  "登录成功",
		"user_id":     req.UserName,
	})
}
func RegisterHandle(c *gin.Context) {
	// var req models.RegisterRequest
	// if err := c.ShouldBind(&req); err != nil {
	// 	c.JSON(http.StatusOK, auth.LoginResponse{
	// 		StatusCode: constant.GateWayParamsErrorCode,
	// 		StatusMsg:  constant.GateWayParamsError,
	// 		Uid:        0,
	// 		Token:      "",
	// 	})
	// 	return
	// }

	// res, err := Client.Register(c.Request.Context(), &auth.RegisterRequest{
	// 	Username: req.UserName,
	// 	Password: req.Password,
	// })

	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"Username": req.UserName,
	// 	}).Warnf("Error when trying to register")
	// 	c.JSON(http.StatusOK, res)
	// 	return
	// }

	// logrus.WithFields(logrus.Fields{
	// 	"Username": req.UserName,
	// 	"Token":    res.Token,
	// 	"UserId":   res.Uid,
	// }).Infof("User %s has registered.", req.UserName)

	// c.JSON(http.StatusOK, res)
	c.JSON(200, gin.H{
		"message": "success sign up",
		"token":   "123456",
	})
}

func CheckLogin(c *gin.Context) {
	session := sessions.Default(c)

	// 从 session 里获取用户信息
	username := session.Get("username")
	userid := session.Get("userid")

	if username == nil || userid == nil {
		// session 无效或未登录
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code": 401,
			"status_msg":  "user not logged in",
		})
		return
	}

	// session 有效，用户已登录
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"status_msg":  "user logged in",
		"user_id":     userid,
		"username":    username,
	})
}
