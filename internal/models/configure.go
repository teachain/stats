package models

import (
	"time"
	"xorm.io/xorm"
)

type Configure struct {
	ID        int64     `xorm:"pk autoincr 'id'"`
	Name      string    `xorm:"varchar(50) notnull 'name' comment('参数名称')" json:"name"`
	Value     string    `xorm:"varchar(255) notnull 'value' comment('参数值')" json:"value"`
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
	UpdatedAt time.Time `xorm:"TIMESTAMP null 'updated_at' comment('更新时间')"  json:"updated_at"`
}

func SyncConfigure(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(Configure))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(Configure))
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Configure) Exist(db *xorm.Engine, name string) bool {
	has, err := db.Table(c.TableName()).Where("name = ?", name).Exist()
	if err == nil && has {
		return true
	}
	return false
}
func (c *Configure) TableName() string {
	return "configure"
}
func (c *Configure) UpdateValue(db *xorm.Engine, name string, value string) error {
	sql := "update `configure` set value=?,updated_at=now() where name=?"
	_, err := db.Exec(sql, value, name)
	return err
}
func (c *Configure) Save(db *xorm.Engine) (int64, error) {
	return db.Insert(c)
}
func (c *Configure) First(db *xorm.Engine) (bool, error) {
	return db.Limit(1, 0).OrderBy("id ASC").Get(c)
}
