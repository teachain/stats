package models

import (
	"xorm.io/xorm"
)

func SyncTableStruct(db *xorm.Engine) error {
	err := SyncSourceTx(db)
	if err != nil {
		return err
	}
	err = SyncSourceStats(db)
	if err != nil {
		return err
	}
	err = SyncSourceMonth(db)
	if err != nil {
		return err
	}
	err = SyncSourceDay(db)
	if err != nil {
		return err
	}
	err = SyncConfigure(db)
	if err != nil {
		return err
	}
	err = SyncSourceHour(db)
	if err != nil {
		return err
	}
	return nil
}
