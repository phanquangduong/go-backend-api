package main

import (
	_ "go/go-backend-api/cmd/swag/docs"
	"go/go-backend-api/global"
	"go/go-backend-api/internal/initialize"
	"strconv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API Documentation Ecommerce Backend shopGO
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  https://github.com/phanquangduong/go-backend-api

// @contact.name   quangduong
// @contact.url    https://github.com/phanquangduong/go-backend-api
// @contact.email  phanquangduong2002@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8002
// @BasePath  /v1
// @schema http

func main() {
	r := initialize.Run()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Lấy port từ config
	s := global.Config.Server
	port := strconv.Itoa(s.Port)
	r.Run(":" + port)

}
