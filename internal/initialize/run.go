package initialize

import (
	"fmt"
	"go/go-backend-api/global"
	"strconv"

	"go.uber.org/zap"
)

func Run() {
	//	load configuration
	LoadConfig()
	m := global.Config.Mysql
	fmt.Println("Loading configuration nysql", m.Username, m.Password)
	InitLogger()
	global.Logger.Info("Config Log ok!!", zap.String("ok", "success"))
	InitMysql()
	InitMysqlC()
	InitServiceInterface()
	InitRedis()
	InitKafka()
	s := global.Config.Server
	port := strconv.Itoa(s.Port)

	r := InitRouter()
	r.Run(":" + port)

}
