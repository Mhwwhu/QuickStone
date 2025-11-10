package user

import (
	"net/http"

	"QuickStone/src/config"
	"QuickStone/src/constant"
	"QuickStone/src/models/webModels"
	"QuickStone/src/rpc/user"
	grpcutil "QuickStone/src/utils/grpc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var userClient user.UserServiceClient

func init() {
	conn := grpcutil.Connect(config.UserServerName)
	userClient = user.NewUserServiceClient(conn)
}

// func LoginHandle(c *gin.Context) {
// 	var req models.LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{
// 			"status_code": 400,
// 			"status_msg":  "参数错误",
// 		})
// 		return
// 	}

// 	// TODO: 这里用真实数据库验证用户密码
// 	if !checkUserPassword(req.UserName, req.Password) {
// 		c.JSON(401, gin.H{
// 			"status_code": 401,
// 			"status_msg":  "用户名或密码错误",
// 		})
// 		return
// 	}

// 	session := sessions.Default(c)
// 	session.Set("userid", req.UserName)
// 	session.Set("username", req.UserName)
// 	if err := session.Save(); err != nil {
// 		c.JSON(500, gin.H{
// 			"status_code": 500,
// 			"status_msg":  "保存 session 失败",
// 		})
// 		return
// 	}

// 	c.JSON(200, gin.H{
// 		"status_code": 200,
// 		"status_msg":  "登录成功",
// 		"user_id":     req.UserName,
// 	})
// }

// 示例用户验证函数，改成真实 DB 查询即可
// func checkUserPassword(username, password string) bool {

// 	users := map[string]string{
// 		"123@qq.com": "123",
// 		"bob":        "123",
// 	}
// 	pwd, ok := users[username]
// 	return ok && pwd == password
// }

func RegisterHandle(c *gin.Context) {
	var req webModels.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.RegisterResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	res, err := userClient.Register(c.Request.Context(), &user.RegisterRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Info("Error when trying to register")
		c.JSON(http.StatusOK, webModels.RegisterResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username":   req.UserName,
		"Token":      res.Token,
		"UserId":     res.Uid,
		"StatusCode": res.StatusCode,
	}).Infof("User %s has registered.", req.UserName)

	c.JSON(http.StatusOK, webModels.RegisterResponse{
		StandardResponse: webModels.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		UserId: res.Uid,
		Token:  res.Token,
	})

}

func LoginHandle(c *gin.Context) {
	var req webModels.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, webModels.LoginResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: constant.GateWayParamsErrorCode,
				StatusMsg:  constant.GateWayParamsError,
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	res, err := userClient.Login(c.Request.Context(), &user.LoginRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Infof("Error when trying to login: %v", err)
		c.JSON(http.StatusOK, webModels.LoginResponse{
			StandardResponse: webModels.StandardResponse{
				StatusCode: res.StatusCode,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.Uid,
	}).Infof("User %s has login.", req.UserName)

	c.JSON(http.StatusOK, webModels.LoginResponse{
		StandardResponse: webModels.StandardResponse{
			StatusCode: res.StatusCode,
			StatusMsg:  res.StatusMsg,
		},
		UserId: res.Uid,
		Token:  res.Token,
	})
}

func LogoutHandle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "logout success",
	})
}
