package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go/go-backend-api/global"
	"go/go-backend-api/internal/constants"
	"go/go-backend-api/internal/database"
	"go/go-backend-api/internal/models"
	"go/go-backend-api/internal/utils"
	"go/go-backend-api/internal/utils/auth"
	"go/go-backend-api/internal/utils/crypto"
	"go/go-backend-api/internal/utils/random"
	"go/go-backend-api/pkg/response"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type sUserLogin struct {
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

// ------ TWO FACTOR AUTHENTICATION -----

func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (codeResult int, rs bool, err error) {
	return 200, true, nil
}
func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *models.SetupTwoFactorAuthInput) (codeResult int, err error) {
	// 1. Check isTwoFactorEnabled -> true -> return
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, in.UserId)
	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}

	if isTwoFactorAuth > 0 {
		return response.ErrCodeTwoFactorAuthSetupFailed, fmt.Errorf("two factor authentication is already enabled")
	}

	// 2. Create new type Auth
	err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.PreGoAccUserTwoFactorTwoFactorAuthTypeEMAIL,
		TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
	})

	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}

	// 3. Send otp to in.TwoFactorEmail
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	// otpNew := random.GenerateSixDigitOtp()
	go global.Rdb.Set(ctx, keyUserTwoFactor, "123456", time.Duration(constants.TIME_2FA_OTP)*time.Minute).Err()

	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *models.TwoFactorVerificationAuthInput) (codeResult int, err error) {
	// 1. Check isTwoFactorEnabled
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, in.UserId)
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	if isTwoFactorAuth > 0 {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("two factor authentication is alrealy enabled")
	}

	// 2. Check otp in redis available
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))

	otpVerifyAuth, err := global.Rdb.Get(ctx, keyUserTwoFactor).Result()
	if err == redis.Nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("key %s does not exists", keyUserTwoFactor)
	}

	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	// 3. Check otp
	if otpVerifyAuth != in.TwoFactorCode {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("OTP not match")
	}

	// 4. update status
	err = s.r.UpdateTwoFactorStatus(ctx, database.UpdateTwoFactorStatusParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.PreGoAccUserTwoFactorTwoFactorAuthTypeEMAIL,
	})

	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	err = s.r.UpdateIsTwoFactorEnabled(ctx, int32(in.UserId))

	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}
	// 5. Remove OTP
	_, err = global.Rdb.Del(ctx, keyUserTwoFactor).Result()
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	return 200, nil
}

// ----- END TWO FACTOR AUTHENTICATION -----

// Implement the IUserLogin interface here
func (s *sUserLogin) Login(ctx context.Context, in *models.LoginInput) (codeResult int, out models.LoginOutput, err error) {
	// 1. logic login
	userBase, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// 2. check password?
	if !crypto.MatchingPassword(userBase.UserPassword, in.UserPassword, userBase.UserSalt) {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("does not match password")
	}
	// 3. check two-factor authentication
	isTwoFactorEnabled, err := s.r.IsTwoFactorEnabled(ctx, uint32(userBase.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}

	if isTwoFactorEnabled > 0 {
		// send otp to in.TwoFactorEmail
		keyUserLoginTwoFactor := crypto.GetHash("2fa:otp" + strconv.Itoa(int(userBase.UserID)))
		err = global.Rdb.SetEx(ctx, keyUserLoginTwoFactor, "111111", time.Duration(constants.TIME_2FA_OTP)*time.Minute).Err()
		if err != nil {
			return response.ErrCodeAuthFailed, out, fmt.Errorf("set otp redis failed")
		}

		// send otp via two factorEmail
		// get email 2fa
		infoUserTwoFactor, err := s.r.GetTwoFactorMethodByIDAndType(ctx, database.GetTwoFactorMethodByIDAndTypeParams{
			UserID:            uint32(userBase.UserID),
			TwoFactorAuthType: database.PreGoAccUserTwoFactorTwoFactorAuthTypeEMAIL,
		})

		if err != nil {
			return response.ErrCodeAuthFailed, out, fmt.Errorf("get two factor method failed")
		}
		// go sendto mail
		log.Println("send OTP 2FA to Email::", infoUserTwoFactor.TwoFactorEmail)
		// go send.SendTextEmailOtp([]string{infoUserTwoFactor.TwoFactorEmail.String}, consts.HOST_EMAIL, "111111")

		out.Message = "send OTP 2FA to Email, pls get OTP by email...."
		return response.ErrCodeSuccess, out, nil
	}

	// 4. update password time
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp: sql.NullString{String: "127.0.0.1", Valid: true},
		UserAccount: in.UserAccount,
	})

	// 5. Create UUID User
	subToken := utils.GenerateCliTokenUUID(int(userBase.UserID))
	log.Println("subtoken:", subToken)
	// 6. get user_info table
	infoUser, err := s.r.GetUser(ctx, uint64(userBase.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// convert to json
	infoUserJson, err := json.Marshal(infoUser)
	if err != nil {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("convert to json failed: %v", err)
	}
	// 7. give infoUserJson to redis with key = subToken
	err = global.Rdb.Set(ctx, subToken, infoUserJson, time.Duration(constants.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// 8. create token
	out.Token, err = auth.CreateToken(subToken)
	if err != nil {
		return
	}
	return response.ErrCodeSuccess, out, nil
}

func (s *sUserLogin) Register(ctx context.Context, in *models.RegisterInput) (codeResult int, err error) {
	// logic
	// 1. hash email
	fmt.Printf("VerifyKey: %s\n", in.VerifyKey)
	fmt.Printf("VerifyType: %d\n", in.VerifyType)
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("hashKey: %s\n", hashKey)

	// 2. check user exists in user base
	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	if userFound > 0 {
		return response.ErrCodeUserHasExists, fmt.Errorf("user has already registered")
	}

	// 3. Create OTP
	userKey := utils.GetUserKey(hashKey) //fmt.Sprintf("u:%s:otp", hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	// util..
	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed::", err)
		return response.ErrInvalidOTP, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("")
	}

	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}
	fmt.Printf("Otp is :::%d\n", otpNew)
	// 5. save OTP in Redis with expiration time
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(constants.TIME_OTP_REGISTER)*time.Minute).Err()

	if err != nil {
		return response.ErrInvalidOTP, err
	}
	// 6. Send OTP
	switch in.VerifyType {
	case constants.EMAIL:
		// err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, constants.HOST_EMAIL, strconv.Itoa(otpNew))
		// if err != nil {
		// 	return response.ErrSendEmailOTP, fmt.Errorf("failed to send email OTP")
		// }
		// 7. save OTP to MYSQL
		result, err := s.r.InsertOrUpdateOTPVerify(ctx, database.InsertOrUpdateOTPVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
			VerifyKey:     in.VerifyKey,
			VerifyKeyHash: hashKey,
		})

		if err != nil {
			return response.ErrSendEmailOTP, err
		}

		// 8. getlasId
		lastIdVerifyUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrSendEmailOTP, err
		}
		log.Println("lastIdVerifyUser", lastIdVerifyUser)
		return response.ErrCodeSuccess, nil
	case constants.MOBILE:
		return response.ErrCodeSuccess, nil
	}

	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyOTP(ctx context.Context, in *models.VerifyInput) (out models.VerifyOtpOutput, err error) {
	// logic
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	// get otp
	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}
	if in.VerifyCode != otpFound {
		// Neu nhu otp sai 3 lan trong vong 1 phut

		return out, fmt.Errorf("OTP not match")
	}

	infoOtp, err := s.r.GetInfoOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}

	// update status verified
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, err
	}

	out.Token = infoOtp.VerifyKeyHash
	out.Message = "success"

	return out, err
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	// token is already verified : user_verify table
	infoOtp, err := s.r.GetInfoOTP(ctx, token)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	// check isVerified ok
	if infoOtp.IsVerified.Int32 == 0 {
		return response.ErrCodeOtpNotExists, fmt.Errorf("user OTP not verified")
	}

	// check token is exists in user_base table

	// update user_base table
	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = infoOtp.VerifyKey
	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HashPassword(password, userSalt)
	// add UserBase to user_base table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}

	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}

	// add user_id to user_info table
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:               uint64(user_id),
		UserAccount:          infoOtp.VerifyKey,
		UserNickname:         sql.NullString{String: infoOtp.VerifyKey, Valid: true},
		UserAvatar:           sql.NullString{String: "", Valid: true},
		UserState:            1,
		UserMobile:           sql.NullString{String: "", Valid: true},
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:         sql.NullTime{Time: time.Time{}, Valid: false},
		UserEmail:            sql.NullString{String: infoOtp.VerifyKey, Valid: true},
		UserIsAuthentication: 1,
	})

	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	user_id, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrCodeOtpNotExists, err
	}

	return int(user_id), nil
}
