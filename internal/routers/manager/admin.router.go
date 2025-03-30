package manager

import "github.com/gin-gonic/gin"

type AdminRouter struct{}

func (pr *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {

	// public router
	adminRouterPublic := Router.Group("admin")
	{
		adminRouterPublic.POST("/login")
	}

}
