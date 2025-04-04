package account

import (
	"go/go-backend-api/internal/models"
	"go/go-backend-api/internal/service"
	"go/go-backend-api/internal/utils/context"
	"go/go-backend-api/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

var TwoFactor = new(cUserTwoFactor)

type cUserTwoFactor struct {
}

// User setup Two Factor Authentication documentation
// @Summary      User setup Two Factor Authentication
// @Description  User setup Two Factor Authentication
// @Tags         account two factor
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization token"
// @Param        payload body models.SetupTwoFactorAuthInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/two_factor/setup [post]
func (c *cUserTwoFactor) SetupTwoFactorAuth(ctx *gin.Context) {
	var params models.SetupTwoFactorAuthInput

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorAuthSetupFailed, "Missing or invalid setup two factorauth parameter")
		return
	}

	// get userId from uuid (token)
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorAuthSetupFailed, "UserId is not invalid")
		return
	}

	log.Println("User Id: ", userId)
	params.UserId = uint32(userId)
	codeResult, err := service.UserLogin().SetupTwoFactorAuth(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorAuthSetupFailed, err.Error())
		return
	}

	response.SuccessResponse(ctx, codeResult, nil)
}
