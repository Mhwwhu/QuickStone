package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/models"
)

// var Client auth.AuthServiceClient
// func init() {
// 	conn := grpcutil.Connect(config.AuthServerName)
// 	Client = auth.NewAuthServiceClient(conn)
// }

func LoginHandle(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"status_msg":  "参数错误",
		})
		return
	}

	// TODO: 这里用真实数据库验证用户密码
	if !checkUserPassword(req.UserName, req.Password) {
		c.JSON(401, gin.H{
			"status_code": 401,
			"status_msg":  "用户名或密码错误",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("userid", req.UserName)
	session.Set("username", req.UserName)
	if err := session.Save(); err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "保存 session 失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"status_msg":  "登录成功",
		"user_id":     req.UserName,
	})
}

// 示例用户验证函数，改成真实 DB 查询即可
func checkUserPassword(username, password string) bool {

	users := map[string]string{
		"123@qq.com": "123",
		"bob":        "123",
	}
	pwd, ok := users[username]
	return ok && pwd == password
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
