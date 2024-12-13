package models

import (
	"time"
	"xorm.io/xorm"
)

type SourceMonth struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	//主题
	Source string `xorm:"varchar(50) notnull 'source' comment('数据来源')" json:"source"`
	// nameServer地址
	Total     int64     `xorm:"'total' comment('总数')"  json:"total"`
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
	UpdatedAt time.Time `xorm:"TIMESTAMP null 'updated_at' comment('更新时间')"  json:"updated_at"`
	//统计的月份
	MonthIn string `xorm:"varchar(7) notnull 'month_in' comment('月份')" json:"month_in"`
}

func SyncSourceMonth(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(SourceMonth))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(SourceMonth))
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SourceMonth) Exist(db *xorm.Session, source string, month string) bool {
	has, err := db.Table(s.TableName()).Where("source = ?", source).Where("month_in = ?", month).Exist()
	if err == nil && has {
		return true
	}
	return false
}
func (s *SourceMonth) UpdateTotal(db *xorm.Session, source string, month string) error {
	sql := "update `source_month` set total=total+1,updated_at=now() where source=? and month_in=?"
	_, err := db.Exec(sql, source, month)
	return err
}
func (s *SourceMonth) TableName() string {
	return "source_month"
}
func (s *SourceMonth) Save(db *xorm.Session) (int64, error) {
	return db.Insert(s)
}
