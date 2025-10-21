package main

import (
	"context"

	"github.com/mhwwhu/QuickStone/src/rpc/auth"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	auth.AuthServiceServer
}

func (a AuthService) Init() {

}

func (a AuthService) Login(ctx context.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	logrus.Infof("Service %s is processing Login", ServerId)
	resp = &auth.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "Login success!",
		Uid:        1,
		Token:      "111",
	}
	return
}

func (a AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	logrus.Infof("Service %s is processing Register", ServerId)
	resp = &auth.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "Register success!",
		Uid:        1,
		Token:      "111",
	}
	return
}
