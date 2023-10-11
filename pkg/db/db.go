package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"im/config"
	"im/pkg/logger"
	"im/pkg/util"
	"os"
	"time"
)

var (
	DB        *gorm.DB
	RedisCli  *redis.Client
	RedisUtil *util.RedisUtil
)

func Init() {
	InitMysql()
	InitRedis()
}

// InitMysql 初始化MySQL
func InitMysql() {
	logger.Logger.Info("init mysql")

	cfg := config.Config.MySQL

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", cfg.Username, cfg.Password, cfg.Address, "mysql")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("init mysql error, open mysql error, err: %v", err))
	}

	dbSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DatabaseName)
	if err = DB.Exec(dbSql).Error; err != nil {
		panic(fmt.Sprintf("init mysql error, create database error, err: %v", err))
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger.Default.LogMode(logger.DBLevel),
		DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(fmt.Sprintf("init mysql error, open %s error, err: %v", cfg.DatabaseName, err))
	}

	var sqlDB *sql.DB
	if sqlDB, err = DB.DB(); err != nil {
		panic(fmt.Sprintf("db error, err: %v", err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	logger.Logger.Info("init mysql ok")
}

// InitRedis 初始化Redis
func InitRedis() {
	logger.Logger.Info("init redis")

	cfg := config.Config.Redis

	RedisCli = redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		DB:       0,
		Password: cfg.Password,
	})

	_, err := RedisCli.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("init redis error, ping error, err: %v", err))
	}

	RedisUtil = util.NewRedisUtil(RedisCli)
	logger.Logger.Info("init redis ok")
}

// InitByTest 初始化数据库配置，仅用在单元测试
func InitByTest() {
	fmt.Println("init db")
	logger.Target = logger.Console

	InitMysql()
	InitRedis()
}
