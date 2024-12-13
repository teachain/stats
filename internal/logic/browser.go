package logic

import (
	"context"
	"github.com/simplechain-org/client/core/types"
	"github.com/simplechain-org/client/ethclient"
	"github.com/simplechain-org/client/log"
	"math/big"
	"sync/atomic"
	"time"
)

type Browser interface {
	OnBlockUpdate(number *big.Int)
	GetName() string
}

type BlockBrowser struct {
	client     *ethclient.Client
	timeout    time.Duration
	stop       chan struct{}
	isStopped  int32
	browsers   map[string]Browser
	wsEndpoint string
}

func (b *BlockBrowser) RegisterBrowser(browser Browser) {
	if browser != nil {
		b.browsers[browser.GetName()] = browser
	}
}

func NewBlockBrowser(wsEndpoint string, timeout time.Duration) (*BlockBrowser, error) {
	client, err := ethclient.DialContext(context.Background(), wsEndpoint)
	if err != nil {
		return nil, err
	}
	b := &BlockBrowser{
		client:     client,
		timeout:    timeout,
		stop:       make(chan struct{}),
		isStopped:  0,
		browsers:   make(map[string]Browser),
		wsEndpoint: wsEndpoint,
	}
	return b, nil
}

func (b *BlockBrowser) onNewBlock(number *big.Int) {
	for _, browser := range b.browsers {
		browser.OnBlockUpdate(big.NewInt(0).Set(number))
	}
}

func (b *BlockBrowser) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()
	ch := make(chan *types.Header, 10)
	subscription, err := b.client.SubscribeNewHead(ctx, ch)
	if err != nil {
		log.Error("SubscribeNewHead", "error", err.Error())
		return err
	}
	go func() {
		defer func() {
			log.Info("BlockBrowser Start exited")
		}()
		for {
			select {
			case header := <-ch:
				log.Info("Receive new block",
					"number", header.Number.String(),
					"at", time.Now().Format(time.DateTime))
				b.onNewBlock(header.Number)
			case err := <-subscription.Err():
				if err != nil {
					log.Error("SubscribeNewHead run", "error", err.Error())
				loop:
					for {
						select {
						case <-b.stop:
							log.Info("BlockBrowser Start subscribeNewHead exits")
							return
						default:
							//连接断开以后，需要重新进行订阅，虽然Client自身有重连操作，但不会自动重新订阅
							subscription, err = b.client.SubscribeNewHead(context.Background(), ch)
							if err != nil {
								log.Error("SubscribeNewHead Resubscribe", "error", err.Error())
								time.Sleep(time.Second * 5)
							} else {
								break loop
							}
						}
					}
				}
			case <-b.stop:
				return
			}
		}
	}()
	return nil
}

func (b *BlockBrowser) Stop() {
	if atomic.CompareAndSwapInt32(&b.isStopped, 0, 1) {
		close(b.stop)
	}
}
