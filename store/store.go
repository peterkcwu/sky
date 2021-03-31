package store

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sky/config"
	"time"
)

//创建store 对象
type Store struct {
	db *gorm.DB
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
	db.AutoMigrate(&Approve{}, &Alarm{}, &Device{}, &AlarmVideo{})
	return &Store{db: db}, nil
}

func (store *Store) GetApproves() ([]Approve, error) {
	var lists []Approve
	err := store.db.Order("id DESC").Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (store *Store) GetApprovesFilter(filter Approve) ([]Approve, error) {
	var lists []Approve
	db := store.db
	if filter.ID != 0 {
		db = db.Where("id = ?", filter.ID)
	}
	if filter.Approver != "" {
		db = db.Where("approver = ?", filter.Approver)
	}
	if filter.Reason != "" {
		db = db.Where("reason = ?", filter.Reason)
	}
	err := db.Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (store *Store) CreateApprove(approve Approve) (Approve, error) {
	err := store.db.Create(&approve).Error
	return approve, err
}

func (store *Store) UpdateApprove(approve Approve) (Approve, error) {
	return approve, store.db.Save(&approve).Error
}

func (store *Store) DeleteApprove(id uint) error {
	return store.db.Unscoped().Where("id = ?", id).Delete(&Approve{}).Error
}

func (store *Store) GetApproveByID(id uint) (Approve, error) {
	var approve Approve
	return approve, store.db.Where("id = ?", id).First(&approve).Error
}

func (store *Store) GetAlarms() ([]Alarm, error) {
	var lists []Alarm
	err := store.db.Order("id DESC").Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (store *Store) CreateAlarm(alarm Alarm) (Alarm, error) {
	err := store.db.Create(&alarm).Error
	return alarm, err
}

func (store *Store) UpdateAlarm(Alarm Alarm) (Alarm, error) {
	return Alarm, store.db.Save(&Alarm).Error
}

func (store *Store) DeleteAlarm(id uint) error {
	return store.db.Unscoped().Where("id = ?", id).Delete(&Alarm{}).Error
}

func (store *Store) GetAlarmByID(id uint) (Alarm, error) {
	var alarmVideo Alarm
	return alarmVideo, store.db.Where("id = ?", id).First(&alarmVideo).Error
}

func (store *Store) CreateAlarmVideo(alarm AlarmVideo) (AlarmVideo, error) {
	err := store.db.Create(&alarm).Error
	return alarm, err
}

func (store *Store) UpdateAlarmVideo(Alarm AlarmVideo) (AlarmVideo, error) {
	return Alarm, store.db.Save(&Alarm).Error
}

func (store *Store) DeleteAlarmVideo(id uint) error {
	return store.db.Unscoped().Where("id = ?", id).Delete(&AlarmVideo{}).Error
}

func (store *Store) GetAlarmVideoByID(id uint) (AlarmVideo, error) {
	var alarm AlarmVideo
	return alarm, store.db.Where("id = ?", id).First(&alarm).Error
}
