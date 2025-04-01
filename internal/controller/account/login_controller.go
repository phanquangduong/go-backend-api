package account

import (
	"go/go-backend-api/internal/service"
	"go/go-backend-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// management controller login user
var Login = new(cUserLogin)

type cUserLogin struct {
}

func (c *cUserLogin) Login(ctx *gin.Context) {
	// Implement logic for login
	err := service.UserLogin().Login(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
