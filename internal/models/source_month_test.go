package models

import (
	"fmt"
	"testing"
	"time"
	"xorm.io/xorm"
)

func TestSyncSourceMonth(t *testing.T) {
	dsn := "root:qq851126@tcp(127.0.0.1)/chain_stats?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		t.Error(err.Error())
		return
	}
	db := d.NewSession()
	defer db.Close()
	err = SyncSourceMonth(d)
	if err != nil {
		t.Error(err.Error())
		return
	}

	stats := &SourceMonth{}

	if stats.Exist(db, "baoquan", CurMonth()) {
		err := stats.UpdateTotal(db, "baoquan", CurMonth())
		if err != nil {
			t.Error(err.Error())
			return
		}
	} else {
		stats.Source = "baoquan"
		stats.Total = 1
		stats.CreatedAt = time.Now()
		stats.UpdatedAt = time.Now()
		stats.MonthIn = CurMonth()
		n, err := stats.Save(db)
		if err != nil {
			t.Error(err.Error())
			return
		}
		fmt.Println(n)
	}
}
