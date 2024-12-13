package models

import (
	"time"
	"xorm.io/xorm"
)

type SourceDay struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	//主题
	Source string `xorm:"varchar(50) notnull 'source' comment('数据来源')" json:"source"`
	// nameServer地址
	Total     int64     `xorm:"'total' comment('总数')"  json:"total"`
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
	UpdatedAt time.Time `xorm:"TIMESTAMP null 'updated_at' comment('更新时间')"  json:"updated_at"`
	//统计的日期
	DayOn string `xorm:"varchar(10) notnull 'day_on' comment('具体的日期')" json:"day_on"`
}

func SyncSourceDay(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(SourceDay))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(SourceDay))
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SourceDay) Exist(db *xorm.Session, source string, day string) bool {
	has, err := db.Table(s.TableName()).Where("source = ?", source).Where("day_on = ?", day).Exist()
	if err == nil && has {
		return true
	}
	return false
}
func (s *SourceDay) UpdateTotal(db *xorm.Session, source string, day string) error {
	sql := "update `source_day` set total=total+1,updated_at=now() where source=? and day_on=?"
	_, err := db.Exec(sql, source, day)
	return err
}
func (s *SourceDay) TableName() string {
	return "source_day"
}
func (s *SourceDay) Save(db *xorm.Session) (int64, error) {
	return db.Insert(s)
}
