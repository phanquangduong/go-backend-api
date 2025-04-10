package user

import (
	"go/go-backend-api/internal/controller/account"
	"go/go-backend-api/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router
	// this is non-dependency
	// ur := repo.NewUserRepository()
	// us := service.NewUserService(ur)
	// userHandlerNonDependency := controller.NewUserController(us)

	// WIRE go
	// Dependency Injection (DI)
	// userController, _ := wire.InitUserRouterHandler()

	userRouterPublic := Router.Group("/user")
	{
		userRouterPublic.POST("/register", account.Login.Register) // register -> YES -> NO
		userRouterPublic.POST("/login", account.Login.Login)
		userRouterPublic.POST("/verify_account", account.Login.VerifyOTP)
		userRouterPublic.POST("/update_password_registeer", account.Login.UpdatePasswordRegister)
	}

	// private router
	userRouterPrivate := Router.Group("/user")
	userRouterPrivate.Use(middlewares.AuthMiddlewware())

	// userRouterPrivate.User(limiter())
	// userRouterPrivate.Use(Auth())
	// userRouterPrivate.Use(Permission())
	{
		userRouterPrivate.POST("/two_factor/setup", account.TwoFactor.SetupTwoFactorAuth)
		userRouterPrivate.POST("/two_factor/verify", account.TwoFactor.VerifyTwoFactorAuth)
	}

}
