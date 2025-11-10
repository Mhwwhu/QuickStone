package main

import (
	"context"

	"QuickStone/src/constant"
	"QuickStone/src/models/dbModels"
	"QuickStone/src/rpc/user"
	"QuickStone/src/storage/database"
	"QuickStone/src/utils/jwt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	user.UserServiceServer
}

func (a UserService) Init() {

}

func (a UserService) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	logrus.Infof("Service %s is processing Login", ServerId)

	resp = &user.LoginResponse{}
	userModel := dbModels.User{
		Name: req.Username,
	}

	// TODO: 先从缓存中读取用户的token

	result := database.Client.Where("name = ?", req.Username).WithContext(ctx).Find(&userModel)
	if result.Error != nil {
		resp = &user.LoginResponse{
			StatusCode: constant.DatabaseErrorCode,
		}
		return
	}

	if result.RowsAffected == 0 {
		resp = &user.LoginResponse{
			StatusCode: constant.UserExistsErrorCode,
			StatusMsg:  constant.UserNotExistsError,
		}
		return
	}

	if !checkPasswordHash(req.Password, userModel.Password) {
		resp = &user.LoginResponse{
			StatusCode: constant.LoginFailErrorCode,
			StatusMsg:  constant.LoginFailError,
		}
		return
	}

	token := jwt.GetToken(userModel.Id, userModel.Name)

	resp = &user.LoginResponse{
		StatusCode: 0,
		Uid:        userModel.Id,
		Token:      token,
	}

	return
}

func (a UserService) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	logrus.Infof("Service %s is processing Register", ServerId)
	resp = &user.RegisterResponse{}
	var userModel dbModels.User
	result := database.Client.WithContext(ctx).Limit(1).Where("name = ?", req.Username).Find(&userModel)
	if result.RowsAffected != 0 {
		logrus.Infof("User %s already exists.", req.Username)
		resp = &user.RegisterResponse{
			StatusCode: constant.UserExistsErrorCode,
			StatusMsg:  constant.UserExistsError,
			Uid:        0,
		}
		return
	}

	var hashedPassword string
	if hashedPassword, err = hashPassword(req.Password); err != nil {
		logrus.Errorf("hashPassword failed: %v", err)
		resp = &user.RegisterResponse{
			StatusCode: constant.InternalErrorCode,
			Uid:        0,
		}
		return
	}

	userModel.Name = req.Username
	userModel.Password = hashedPassword
	result = database.Client.WithContext(ctx).Create(&userModel)
	if result.Error != nil {
		logrus.Errorf("Database operation failed: %v", result.Error)
		resp = &user.RegisterResponse{
			StatusCode: constant.DatabaseErrorCode,
			Uid:        0,
		}
	}

	resp.Token = jwt.GetToken(userModel.Id, userModel.Name)
	resp.Uid = userModel.Id

	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
