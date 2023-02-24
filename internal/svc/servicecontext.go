package svc

import (
	"context"
	"douyin-zero/internal/config"
	"douyin-zero/internal/dal/query"
	"douyin-zero/internal/middleware"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	Auth        rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := InitRedisClient(c)

	svcCtx := &ServiceContext{
		Config:      c,
		RedisClient: redisClient,
		DB:          InitDB(c),
		Auth:        middleware.NewAuthMiddleware(redisClient).Handle,
	}

	SyncTask(svcCtx)

	return svcCtx
}

func InitRedisClient(c config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("redis connect ping failed, err:%v", err))
	}
	return redisClient
}

func InitDB(c config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.MySQL.DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("mysql connect failed, err:%v", err))
	}
	query.SetDefault(db.Debug())
	return db
}
