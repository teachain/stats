package models

import (
	"time"
	"xorm.io/xorm"
)

type SourceHour struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	//主题
	Source string `xorm:"varchar(50) notnull 'source' comment('数据来源')" json:"source"`
	// nameServer地址
	Total     int64     `xorm:"'total' comment('总数')"  json:"total"`
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
	UpdatedAt time.Time `xorm:"TIMESTAMP null 'updated_at' comment('更新时间')"  json:"updated_at"`
	//统计的日期
	HourAt string `xorm:"varchar(20) notnull 'hour_at' comment('小时')" json:"hour_at"`
}

func SyncSourceHour(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(SourceHour))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(SourceHour))
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SourceHour) Exist(db *xorm.Session, source string, hour string) bool {
	has, err := db.Table(s.TableName()).Where("source = ?", source).Where("hour_at = ?", hour).Exist()
	if err == nil && has {
		return true
	}
	return false
}
func (s *SourceHour) UpdateTotal(db *xorm.Session, source string, hour string) error {
	sql := "update `source_hour` set total=total+1,updated_at=now() where source=? and hour_at=?"
	_, err := db.Exec(sql, source, hour)
	return err
}
func (s *SourceHour) TableName() string {
	return "source_hour"
}
func (s *SourceHour) Save(db *xorm.Session) (int64, error) {
	return db.Insert(s)
}
