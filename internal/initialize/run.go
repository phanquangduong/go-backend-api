package initialize

import (
	"fmt"
	"go/go-backend-api/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	// Load configuration
	LoadConfig()
	m := global.Config.Mysql
	fmt.Println("Loading configuration mysql", m.Username, m.Password)
	InitLogger()
	global.Logger.Info("Config Log ok!!", zap.String("ok", "success"))
	InitMysql()
	InitMysqlC()
	InitServiceInterface()
	InitRedis()
	InitKafka()

	r := InitRouter()
	return r
}
