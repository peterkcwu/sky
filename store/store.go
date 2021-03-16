package store

import (
	"fmt"
	"github.com/new_web/config"
	"time"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//创建store 对象
type Store struct {
	db          *gorm.DB
	redisClient *redis.Client
}

//创建新的数据库
func NewStore(conf *config.Config) (*Store, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local&timeout=10s",
		conf.DB.User, conf.DB.Pass, conf.DB.Host, conf.DB.Port, conf.DB.Database))
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(60 * time.Second)
	db.AutoMigrate(&Approve{})
	return &Store{db: db}, nil
}
