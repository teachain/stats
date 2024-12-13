package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simplechain-org/client/ethclient"
	"github.com/simplechain-org/client/log"
	"github.com/teachain/stats/internal/models"
	"github.com/teachain/stats/internal/types"
	"math/big"
	"strings"
	"sync/atomic"
	"time"
	"xorm.io/xorm"
)

var one *big.Int = big.NewInt(1)

var LastBlock string = "last_block"

type Worker struct {
	client        *ethclient.Client
	db            *xorm.Engine
	handleNumbers chan *big.Int
	stop          chan struct{}
	isStopped     int32
	startBlock    *big.Int
	maxBlock      *big.Int
}

func NewWorker(client *ethclient.Client, db *xorm.Engine, startBlock *big.Int, maxBlock *big.Int) *Worker {
	return &Worker{
		client:        client,
		db:            db,
		handleNumbers: make(chan *big.Int, 100),
		stop:          make(chan struct{}),
		isStopped:     0,
		startBlock:    big.NewInt(0).Set(startBlock),
		maxBlock:      big.NewInt(0).Set(maxBlock),
	}
}
func (w *Worker) Start() {
	w.handleBlock()
	w.loop()
}

func (w *Worker) handleBlock() {
	configure := &models.Configure{}
	go func() {
		for {
			select {
			case number := <-w.handleNumbers:
				err := w.scanBlock(number)
				if err != nil {
					fmt.Println(err)
				} else {
					//更新配置参数
					err := configure.UpdateValue(w.db, LastBlock, number.String())
					if err != nil {
						fmt.Println(err)
					}
				}
			case <-w.stop:
				return
			}
		}
	}()
}
func (w *Worker) Stop() {
	if atomic.CompareAndSwapInt32(&w.isStopped, 0, 1) {
		close(w.stop)
	}
}

func (w *Worker) OnBlockUpdate(number *big.Int) {
	if w.maxBlock.Cmp(number) < 0 {
		w.maxBlock.Set(number)
	}
}
func (w *Worker) GetName() string {
	return "worker"
}

func (w *Worker) scanBlock(number *big.Int) error {
	block, err := w.client.BlockByNumber(context.Background(), number)
	if err != nil {
		return err
	}
	transactions := block.Transactions()
	for _, tx := range transactions {
		payload := tx.Data()
		if len(payload) > 2 {
			onChainRequest := new(types.OnChainRequest)
			err := json.Unmarshal(payload, onChainRequest)
			if err != nil {
				log.Warn("OnChainRequest Unmarshal", "error", err.Error())
				continue
			}
			if len(onChainRequest.Source) > 0 {
				err := w.handleSource(onChainRequest.Source, block.Time(), tx.Hash().String())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (w *Worker) handleSource(sources []string, handledAt uint64, txHash string) error {
	session := w.db.NewSession()
	defer func() {
		err := session.Close()
		if err != nil {
			log.Error("session close", "error", err.Error())
		}
	}()

	// add Begin() before any action
	if err := session.Begin(); err != nil {
		return err
	}

	for _, source := range sources {
		source = strings.TrimSpace(source)
		if len(source) == 0 {
			continue
		}
		//交易记录
		statsTx := &models.SourceTx{
			Source:    source,
			TxHash:    txHash,
			CreatedAt: time.Unix(int64(handledAt), 0),
		}
		_, err := statsTx.Save(session)
		if err != nil {
			return err
		}
		//总数记录
		sourceStats := &models.SourceStats{}
		if sourceStats.Exist(session, source) {
			err := sourceStats.UpdateTotal(session, source)
			if err != nil {
				return err
			}
		} else {
			sourceStats.Source = source
			sourceStats.CreatedAt = time.Now()
			sourceStats.UpdatedAt = time.Now()
			sourceStats.Total = 1
			_, err := sourceStats.Save(session)
			if err != nil {
				return err
			}
		}
		//按月统计
		sourceMonth := &models.SourceMonth{}
		if sourceMonth.Exist(session, source, models.Month(handledAt)) {
			err := sourceMonth.UpdateTotal(session, source, models.Month(handledAt))
			if err != nil {
				return err
			}
		} else {
			sourceMonth.Source = source
			sourceMonth.Total = 1
			sourceMonth.CreatedAt = time.Now()
			sourceMonth.UpdatedAt = time.Now()
			sourceMonth.MonthIn = models.Month(handledAt)
			_, err := sourceMonth.Save(session)
			if err != nil {
				return err
			}
		}
		//按日统计
		sourceDay := &models.SourceDay{}
		if sourceDay.Exist(session, source, models.Day(handledAt)) {
			err := sourceDay.UpdateTotal(session, source, models.Day(handledAt))
			if err != nil {
				return err
			}
		} else {
			sourceDay.Source = source
			sourceDay.Total = 1
			sourceDay.CreatedAt = time.Now()
			sourceDay.UpdatedAt = time.Now()
			sourceDay.DayOn = models.Day(handledAt)
			_, err := sourceDay.Save(session)
			if err != nil {
				return err
			}
		}
		//按小时统计
		sourceHour := &models.SourceHour{}
		if sourceHour.Exist(session, source, models.Hour(handledAt)) {
			err := sourceHour.UpdateTotal(session, source, models.Hour(handledAt))
			if err != nil {
				return err
			}
		} else {
			sourceHour.Source = source
			sourceHour.Total = 1
			sourceHour.CreatedAt = time.Now()
			sourceHour.UpdatedAt = time.Now()
			sourceHour.HourAt = models.Hour(handledAt)
			_, err := sourceHour.Save(session)
			if err != nil {
				return err
			}
		}
	}
	// add Commit() after all actions
	return session.Commit()
}
func (w *Worker) loop() {
	go func() {
		for {
			select {
			case <-w.stop:
				return
			default:
			}
			if w.startBlock.Cmp(w.maxBlock) < 0 {
				for w.startBlock.Cmp(w.maxBlock) < 0 {
					w.handleNumbers <- big.NewInt(0).Set(w.startBlock)
					w.startBlock.Add(w.startBlock, one)
				}
			} else {
				time.Sleep(time.Second * 5)
			}
		}
	}()
}
