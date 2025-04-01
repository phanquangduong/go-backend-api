package global

import (
	"database/sql"
	"go/go-backend-api/pkg/logger"
	"go/go-backend-api/pkg/setting"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config setting.Config
	Logger *logger.LoggerZap
	Mdb    *gorm.DB
	Rdb    *redis.Client
	Mdbc   *sql.DB
)

/*
Config
Redis
Mysql
*/
