package controller

import (
	"fmt"
	"go/go-backend-api/internal/service"
	"go/go-backend-api/internal/vo"
	"go/go-backend-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var params vo.UserRegistrationRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid)
		return
	}

	fmt.Printf("Email params: %s", params.Email)
	result := uc.userService.Register(params.Email, params.Purpose)
	response.SuccessResponse(c, result, nil)
}

// func (uc *UserController) GetUserByID(c *gin.Context) {
// 	userInfo := uc.userService.GetUserInfo()
// 	response.SuccessResponse(c, response.ErrCodeSuccess, userInfo)
// }

// Interface version
