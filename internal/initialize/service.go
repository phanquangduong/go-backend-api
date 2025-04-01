package initialize

import (
	"go/go-backend-api/global"
	"go/go-backend-api/internal/database"
	"go/go-backend-api/internal/service"
	"go/go-backend-api/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)
	// User service interface
	service.InitUserLogin(impl.NewUserLoginImpl(queries))
}
