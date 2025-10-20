package main

import (
	"context"

	"github.com/mhwwhu/QuickStone/src/rpc/auth"
)

type AuthService struct {
	auth.AuthServiceServer
}

func (a AuthService) Init() {

}

func (a AuthService) LoginService(ctx context.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	resp = &auth.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "Login success!",
		Uid:        1,
		Token:      "111",
	}
	return
}

func (a AuthService) RegisterService(ctx context.Context, req *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	resp = &auth.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "Login success!",
		Uid:        1,
		Token:      "111",
	}
	return
}
