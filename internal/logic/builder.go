package logic

import (
	"context"
	"github.com/simplechain-org/client/ethclient"
	"github.com/teachain/stats/internal/config"
	"github.com/teachain/stats/internal/models"
	"math/big"
	"time"
	"xorm.io/xorm"
)

type Builder struct {
	config  *config.Config
	worker  *Worker
	browser *BlockBrowser
}

func NewBuilder(c *config.Config) (*Builder, error) {
	client, err := ethclient.DialContext(context.Background(), c.NodeURL)
	dsn := config.DataSource(c.DB)
	db, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = models.SyncTableStruct(db)
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	maxBlock := big.NewInt(0).SetUint64(blockNumber)
	configure := &models.Configure{}
	//默认从第一个区块开始
	startBlock := big.NewInt(1)
	if configure.Exist(db, LastBlock) {
		ok, err := configure.First(db)
		if err == nil && ok {
			//记录中表示的是已经遍历过的区块
			startBlock.SetString(configure.Value, 10)
			//从记录的区块高度的下一个区块开始
			startBlock.Add(startBlock, big.NewInt(1))
		}
	} else {
		configure.Name = LastBlock
		configure.Value = "0"
		configure.CreatedAt = time.Now()
		configure.UpdatedAt = time.Now()
		_, err := configure.Save(db)
		if err != nil {
			return nil, err
		}
	}
	worker := NewWorker(client, db, startBlock, maxBlock)
	browser, err := NewBlockBrowser(c.NodeURL, time.Second*15)
	if err != nil {
		return nil, err
	}
	browser.RegisterBrowser(worker)
	builder := &Builder{
		config:  c,
		worker:  worker,
		browser: browser,
	}
	return builder, nil
}

func (b *Builder) Start() error {
	b.worker.Start()
	err := b.browser.Start()
	if err != nil {
		return err
	}
	return nil
}
func (b *Builder) Stop() {
	b.worker.Stop()
	b.browser.Stop()
}
