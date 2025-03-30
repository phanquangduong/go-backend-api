//go:build wireinject

package wire

import (
	"go/go-backend-api/internal/controller"
	"go/go-backend-api/internal/repo"
	"go/go-backend-api/internal/service"

	"github.com/google/wire"
)

func InitUserRouterHandler() (*controller.UserController, error) {
	wire.Build(
		repo.NewUserRepository,
		repo.NewUserAuthRepository,
		service.NewUserService,
		controller.NewUserController,
	)

	return new(controller.UserController), nil
}
