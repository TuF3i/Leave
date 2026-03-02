package dao

import (
	"leave/core/pkg/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dao struct {
	pgdb      *gorm.DB
	rdb       *redis.Client
	miniRedis *miniredis.Miniredis
	conf      *config.Config
}

func GetDao(conf *config.Config) (*Dao, error) {
	d := &Dao{conf: conf}

	err := d.initPgSQL()
	if err != nil {
		return nil, err
	}

	err = d.initMiniRedis()
	if err != nil {
		return nil, err
	}

	return d, err
}
