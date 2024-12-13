package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
	"xorm.io/xorm"
)

func TestSyncSourceStats(t *testing.T) {
	dsn := "root:qq851126@tcp(127.0.0.1)/chain_stats?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		t.Error(err.Error())
		return
	}
	db := d.NewSession()
	defer db.Close()
	err = SyncSourceStats(d)
	if err != nil {
		t.Error(err.Error())
		return
	}
	stats := &SourceStats{}
	if stats.Exist(db, "baoquan") {
		err := stats.UpdateTotal(db, "baoquan")
		if err != nil {
			t.Error(err.Error())
			return
		}
	} else {
		stats.Source = "baoquan"
		stats.CreatedAt = time.Now()
		stats.UpdatedAt = time.Now()
		stats.Total = 1
		n, err := stats.Save(db)
		if err != nil {
			t.Error(err.Error())
			return
		}
		fmt.Println(n)
	}

}