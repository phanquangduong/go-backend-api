package repo

import (
	"fmt"
	"go/go-backend-api/global"
	"time"
)

type IUserAuthRepository interface {
	AddOTP(email string, otp int, expirationTime int64) error
}

type userAuthRepository struct{}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}

// AddOTP implements IUserAuthRepository.
func (u *userAuthRepository) AddOTP(email string, otp int, expirationTime int64) error {
	key := fmt.Sprintf("usr:%s:otp", email) // usr:email:otp

	return global.Rdb.SetEx(ctx, key, otp, time.Duration(expirationTime)).Err()
}
