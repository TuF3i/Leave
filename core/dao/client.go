package dao

import (
	"context"
	"fmt"
	"leave/core/models"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (r *Dao[T]) initPgSQL() error {
	// 创建连接URL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s "+
			"sslmode=disable "+
			"connect_timeout=5 "+
			"target_session_attrs=read-write", // 确保连接主库写数据

		r.conf.PgSQL.Addr,
		r.conf.PgSQL.Port,
		r.conf.PgSQL.User,
		r.conf.PgSQL.Passwd,
		r.conf.PgSQL.DBName,
	)

	// 连接pgsql
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return fmt.Errorf("connect to pgsql error: %v", err.Error())
	}

	// 数据库迁移
	err = db.AutoMigrate(
		&models.Comment{},
		&models.LeaveArticle{},
		&models.LeaveUser{},
		&models.Reply{},
		&models.LeaveMsg{},
		&models.Tag{},
	)

	if err != nil {
		return fmt.Errorf("auto migrate error: %v", err.Error())
	}

	r.pgdb = db
	return nil
}

func (r *Dao[T]) initMiniRedis() error {
	s, err := miniredis.Run()
	if err != nil {
		return fmt.Errorf("start miniredis error: %v", err.Error())
	}

	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("connect to miniredis error: %v", err.Error())
	}

	r.miniRedis = s
	r.rdb = rdb

	return nil
}
