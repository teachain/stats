package models

import (
	"time"
	"xorm.io/xorm"
)

type SourceTx struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	//主题
	Source string `xorm:"index varchar(50) notnull 'source' comment('数据来源')" json:"source"`
	//该主题的相关的交易哈希
	TxHash string `xorm:"index varchar(66) notnull 'tx_hash' comment('交易哈希')" json:"tx_hash"`
	//区块时间
	CreatedAt time.Time `xorm:"TIMESTAMP null 'created_at' comment('创建时间')"  json:"created_at"`
}

func (s *SourceTx) TableName() string {
	return "source_tx"
}
func (s *SourceTx) Save(db *xorm.Session) (int64, error) {
	return db.Insert(s)
}
func SyncSourceTx(db *xorm.Engine) error {
	ok, err := db.IsTableExist(new(SourceTx))
	if err != nil {
		return err
	}
	if !ok {
		err := db.Sync2(new(SourceTx))
		if err != nil {
			return err
		}
	}
	return nil
}
