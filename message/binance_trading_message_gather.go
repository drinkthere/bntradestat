package message

import (
	"bntradestat/config"
	"bntradestat/container"
	"bntradestat/context"
	"bntradestat/utils"
	"bntradestat/utils/logger"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"strconv"
	"time"
)

func StartGatherBinanceAggTradeTrade(cfg *config.Config, globalContext *context.GlobalContext, binanceFuturesAggTradeChan chan *binanceFutures.WsAggTradeEvent) {

	go func() {
		defer func() {
			logger.Warn("[AggTradeGather] Binance Futures AggTrade Gather Exited.")
		}()
		startTime := time.Now()
		for {
			trade := <-binanceFuturesAggTradeChan
			instID := trade.Symbol
			if !utils.InArray(instID, cfg.InstIDs) {
				continue
			}

			tradeID := trade.AggregateTradeID
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				logger.Error("[AggTradeGather] parse price %s error %s", trade.Price, err.Error())
				continue
			}

			quantity, err := strconv.ParseFloat(trade.Quantity, 64)
			if err != nil {
				logger.Error("[AggTradeGather] parse quantity %s error %s", trade.Quantity, err.Error())
				continue
			}
			if price < cfg.MinAccuracy || quantity < cfg.MinAccuracy {
				logger.Error("[AggTradeGather] price %f or quantity %f too small", price, quantity)
				continue
			}

			volume := price * quantity

			// 更新trade消息
			globalContext.BinanceFuturesTradingComposite.UpdateTradingInfo(
				instID, tradeID,
				container.VolumeWrapper{Volume: volume, TradeTimeMs: trade.TradeTime},
				cfg.RollingPeriodSeconds)

			if time.Since(startTime) > time.Duration(cfg.RollingPeriodSeconds)*time.Second {
				// 获取trading信息
				tradingInfo := globalContext.BinanceFuturesTradingComposite.StatTradingInfo()
				//logger.Info("[TradingInfo] volume=%f, times=%d", tradingInfo.Volume, tradingInfo.Times)
				if tradingInfo.Times >= cfg.MinTradingTimes || tradingInfo.Volume >= cfg.MinTradingVolume {
					oldStatus := globalContext.BinanceFuturesTradingComposite.GetTradingStatus()
					if oldStatus != config.TradingStatusOK {
						// 状态变为正常交易
						globalContext.BinanceFuturesTradingComposite.UpdateTradingStatus(config.TradingStatusOK)

						// 发送交易状态
						globalContext.TradingStatusCh <- string(config.TradingStatusOK)

						logger.Info("[AggTradeGather] trading status changed from %s to %s, volume=%f, times=%d, vthres=%f, ttrhes=%d",
							oldStatus, config.TradingStatusOK, tradingInfo.Volume, tradingInfo.Times, cfg.MinTradingVolume, cfg.MinTradingTimes)
					}
				} else {
					oldStatus := globalContext.BinanceFuturesTradingComposite.GetTradingStatus()
					if oldStatus != config.TradingStatusStop {
						// 状态变为停止交易
						globalContext.BinanceFuturesTradingComposite.UpdateTradingStatus(config.TradingStatusStop)

						// 发送交易状态
						globalContext.TradingStatusCh <- string(config.TradingStatusStop)

						logger.Info("[AggTradeGather] trading status changed from %s to %s, volume=%f, times=%d, vthres=%f, ttrhes=%d",
							oldStatus, config.TradingStatusStop, tradingInfo.Volume, tradingInfo.Times, cfg.MinTradingVolume, cfg.MinTradingTimes)
					}
				}
			}

		}
	}()
	logger.Info("[AggTradeGather] Start Gather Binance Futures AggTrade")
}
