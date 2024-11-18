package message

import (
	"bntradestat/config"
	"bntradestat/utils/logger"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"math/rand"
	"time"
)

func StartBinanceAggTradeWs(cfg *config.Config, binanceFuturesAggTradeChan chan *binanceFutures.WsAggTradeEvent) {
	binanceSpotWs := newBinanceFuturesAggTradeWebSocket(binanceFuturesAggTradeChan)
	for _, source := range cfg.Sources {
		// 循环不同的IP，监听对应的tickers
		binanceSpotWs.startBinanceAggTradeWs(cfg.InstIDs, source)
		time.Sleep(1 * time.Second)
	}
}

type BinanceFuturesAggTradeWebSocket struct {
	tradeChan chan *binanceFutures.WsAggTradeEvent
}

func newBinanceFuturesAggTradeWebSocket(tradeChan chan *binanceFutures.WsAggTradeEvent) *BinanceFuturesAggTradeWebSocket {
	return &BinanceFuturesAggTradeWebSocket{
		tradeChan: tradeChan,
	}
}

func (ws *BinanceFuturesAggTradeWebSocket) startBinanceAggTradeWs(instIDs []string, source config.Source) {
	innerSpot := newInnerBinanceFuturesAggTradeWebSocket(ws.tradeChan, source)
	innerSpot.subscribeAggTrade(instIDs)
	logger.Info("[AggTradeWs] Start Listen Binance Futures AggTrade Event ip:%s private:%t", source.IP, source.Colo)
}

type innerBinanceFuturesAggTradeWebSocket struct {
	source    config.Source
	instIDs   []string
	tradeChan chan *binanceFutures.WsAggTradeEvent
	isStopped bool
	stopChan  chan struct{}
	randGen   *rand.Rand
}

func newInnerBinanceFuturesAggTradeWebSocket(tradeChan chan *binanceFutures.WsAggTradeEvent, source config.Source) *innerBinanceFuturesAggTradeWebSocket {
	return &innerBinanceFuturesAggTradeWebSocket{
		source:    source,
		tradeChan: tradeChan,
		isStopped: true,
		stopChan:  make(chan struct{}),
		randGen:   rand.New(rand.NewSource(2)),
	}
}

func (iws *innerBinanceFuturesAggTradeWebSocket) handleAggTradeEvent(event *binanceFutures.WsAggTradeEvent) {
	if iws.randGen.Int31n(10000) < 2 {
		logger.Info("[BSTickerWebSocket] Binance Spot Event: %+v", event)
	}

	iws.tradeChan <- event
}

func (iws *innerBinanceFuturesAggTradeWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BSTickerWebSocket] Binance Spot Handle Error And Reconnect Ws: %s", err.Error())
	iws.stopChan <- struct{}{}
	iws.isStopped = true
}

func (iws *innerBinanceFuturesAggTradeWebSocket) subscribeAggTrade(instIDs []string) {

	go func() {
		defer func() {
			logger.Warn("[AggTradeWs] Binance Futures AggTrade Listening Exited.")
		}()
		for {
			if !iws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			var stopChan chan struct{}
			var err error
			if iws.source.Colo {
				binanceFutures.UseIntranet = true
			} else {
				binanceFutures.UseIntranet = false
			}
			if iws.source.IP == "" {
				_, stopChan, err = binanceFutures.WsCombinedAggTradeServe(instIDs, iws.handleAggTradeEvent, iws.handleError)
			} else {
				_, stopChan, err = binanceFutures.WsCombinedAggTradeServeWithIP(iws.source.IP, instIDs, iws.handleAggTradeEvent, iws.handleError)
			}
			if err != nil {
				logger.Error("[AggTradeWs] Subscribe Binance AggTrade Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[AggTradeWs] Subscribe Binance AggTrade: %d", len(instIDs))
			iws.stopChan = stopChan
			iws.isStopped = false
		}
	}()
}
