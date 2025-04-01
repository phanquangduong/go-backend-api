package response

const (
	ErrCodeSuccess      = 20001
	ErrCodeParamInvalid = 20003
	ErrInvalidToken     = 30001
	ErrInvalidOTP       = 30002
	ErrSendEmailOTP     = 30003
	// Register code
	ErrCodeUserHasExists = 50001

	// Error login
	ErrCodeOtpNotExists = 60009
)

var msg = map[int]string{
	ErrCodeSuccess:       "success",
	ErrCodeParamInvalid:  "Email is invalid",
	ErrInvalidToken:      "token is invalid",
	ErrInvalidOTP:        "OTP error",
	ErrSendEmailOTP:      "Failed to send email OTP",
	ErrCodeUserHasExists: "User has already registered",
	ErrCodeOtpNotExists:  "OTP exists but not registered",
}
