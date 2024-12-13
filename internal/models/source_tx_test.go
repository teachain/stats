package models

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
	"xorm.io/xorm"
)

func TestSourceTx_Save(t *testing.T) {
	dsn := "root:qq851126@tcp(127.0.0.1)/chain_stats?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = SyncSourceTx(db)
	if err != nil {
		t.Error(err.Error())
		return
	}
	stats := &SourceTx{
		Source:    "baoquan",
		TxHash:    "hello",
		CreatedAt: time.Now(),
	}
	n, err := stats.Save(db)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("n=", n)

}
