package main

import (
	"bntradestat/message"
	binanceFutures "github.com/dictxwang/go-binance/futures"
)

func startTradeMessage() {
	// 监听binance public trade信息并收集整理
	binanceFuturesAggTradeChan := make(chan *binanceFutures.WsAggTradeEvent)
	message.StartBinanceAggTradeWs(&globalConfig, binanceFuturesAggTradeChan)
	message.StartGatherBinanceAggTradeTrade(&globalConfig, &globalContext, binanceFuturesAggTradeChan)
}
