package service

import (
	"fmt"
	"go/go-backend-api/internal/repo"
	"go/go-backend-api/internal/utils/crypto"
	"go/go-backend-api/internal/utils/random"
	"go/go-backend-api/pkg/response"
	"time"
)

type IUserService interface {
	Register(email string, purpose string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
}

func NewUserService(userRepo repo.IUserRepository, userAuthRepo repo.IUserAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

// Register implements IUserService.
func (us *userService) Register(email string, purpose string) int {
	// 0. hashEmail
	hashEmail := crypto.GetHash(email)
	fmt.Printf("hashEmail::%s\n", hashEmail)

	// 5. check OTP is available

	// 6. user spam ...

	// 1. check email exists in db

	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExists
	}

	// 2. new OTP -> ...
	otp := random.GenerateSixDigitOtp()

	if purpose == "TEST_USER" {
		otp = 123456
	}

	fmt.Printf("Otp is ::%d\n", otp)
	// 3. save OTP in Redis with expiration time

	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))

	if err != nil {
		return response.ErrInvalidOTP
	}
	// 4. send Email OTP

	// err = sendto.SendTemplateEmailOtp([]string{email}, "defnotqduong@gmail.com", "otp-auth.html", map[string]interface{}{
	// 	"otp": strconv.Itoa(otp),
	// })

	// if err != nil {
	// 	return response.ErrSendEmailOTP
	// }

	return response.ErrCodeSuccess

}
