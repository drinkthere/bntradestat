package context

import (
	"bntradestat/config"
	"bntradestat/container"
)

type GlobalContext struct {
	BinanceFuturesTradingComposite *container.TradingComposite
	// 传递tradingStatus的channel
	TradingStatusCh chan string
}

func (context *GlobalContext) Init(globalConfig *config.Config) {
	// 初始化trading数据
	context.initTradingComposite(globalConfig)

	context.TradingStatusCh = make(chan string)
}

func (context *GlobalContext) initTradingComposite(globalConfig *config.Config) {
	context.BinanceFuturesTradingComposite = container.NewTradingComposite(globalConfig)
}
