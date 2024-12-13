package models

import (
	"time"
	"xorm.io/xorm"
)

type SourceStats struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	//主题
	Source string `xorm:"varchar(50) notnull 'source' comment('数据来源')" json:"source"`
	// nameServer地址
	Total     int64     `xorm:"'total' comment('总数')"  json:"total"`
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
	UpdatedAt time.Time `xorm:"TIMESTAMP null 'updated_at' comment('更新时间')"  json:"updated_at"`
}

func SyncSourceStats(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(SourceStats))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(SourceStats))
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SourceStats) IsTableEmpty(db *xorm.Engine) (bool, error) {
	return db.IsTableEmpty(new(SourceStats))
}

func (s *SourceStats) Save(db *xorm.Session) (int64, error) {
	return db.Insert(s)
}
func (s *SourceStats) First(db *xorm.Engine) (bool, error) {
	return db.Limit(1, 0).OrderBy("id ASC").Get(s)
}
func (s *SourceStats) Clear(db *xorm.Engine) error {
	_, err := db.Where("1=1").Delete(new(SourceStats))
	return err
}

func (s *SourceStats) Exist(db *xorm.Session, source string) bool {
	has, err := db.Table(s.TableName()).Where("source = ?", source).Exist()
	if err == nil && has {
		return true
	}
	return false
}
func (s *SourceStats) UpdateTotal(db *xorm.Session, source string) error {
	sql := "update `source_stats` set total=total+1,updated_at=now() where source=?"
	_, err := db.Exec(sql, source)
	return err
}
func (s *SourceStats) TableName() string {
	return "source_stats"
}
