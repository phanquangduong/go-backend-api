package service

import (
	"context"
	"go/go-backend-api/internal/models"
)

type (
	// ... interface
	IUserLogin interface {
		Login(ctx context.Context, in *models.LoginInput) (codeResult int, out models.LoginOutput, err error)
		Register(ctx context.Context, in *models.RegisterInput) (codeResult int, err error)
		VerifyOTP(ctx context.Context, in *models.VerifyInput) (out models.VerifyOtpOutput, err error)
		UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error)
	}

	IUserInfo interface {
		GetInfoByUserId(ctx context.Context) error
		GetAllUser(ctx context.Context) error
	}

	IUserAdmin interface {
		RemoveUser(ctx context.Context) error
		FindOneUser(ctx context.Context) error
	}
)

var (
	localUserAdmin IUserAdmin
	localUserInfo  IUserInfo
	localUserLogin IUserLogin
)

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("implement localUserAdmin not found for interface IUserAdmin")
	}
	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("implement localUserInfo not found for interface IUserInfo")
	}
	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}

func UserLogin() IUserLogin {
	if localUserLogin == nil {
		panic("implement localUserLogin not found for interface IUserLogin")
	}
	return localUserLogin
}

func InitUserLogin(i IUserLogin) {
	localUserLogin = i
}
