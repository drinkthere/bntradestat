package main

import (
	"bntradestat/protocol/pb"
	"bntradestat/utils/logger"
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
)

func StartZmq() {
	logger.Info("Start Binance Trading Status ZMQ")
	startBinanceTradingStatusZmq()
}

func startBinanceTradingStatusZmq() {
	go func() {
		defer func() {
			logger.Warn("[BNTradingStatusZmq] Pub Service Listening Exited.")
		}()

		logger.Info("[BNTradingStatusZmq] Start Pub Service.")

		var ctx *zmq.Context
		var pub *zmq.Socket
		var err error
		for {
			ctx, err = zmq.NewContext()
			if err != nil {
				logger.Error("[BNTradingStatusZmq] New Context Error: %s", err.Error())
				continue
			}

			pub, err = ctx.NewSocket(zmq.PUB)
			if err != nil {
				logger.Error("[BNTradingStatusZmq] New Socket Error: %s", err.Error())
				ctx.Term()
				continue
			}

			ipc := globalConfig.TradingStatusZMQIPC
			err = pub.Bind(ipc)
			if err != nil {
				logger.Error("[BNFuturesZmq] Bind to  ZMQ %s Error: %s", ipc, err.Error())
				pub.Close()
				ctx.Term()
				continue
			}

			for {
				select {
				case tradingStatus, ok := <-globalContext.TradingStatusCh:
					if !ok {
						logger.Warn("[BNFuturesTradingStatusZmq] Trading status channel closed.")
						pub.Close()
						ctx.Term()
						return
					}

					md := &pb.TradingStatus{
						Status: tradingStatus,
					}

					data, err := proto.Marshal(md)
					if err != nil {
						logger.Warn("[BNFuturesTradingStatusZmq] Error marshaling trading status %s, error: %+v", tradingStatus, err)
						continue
					}

					_, err = pub.Send(string(data), 0)
					if err != nil {
						logger.Warn("[BNFuturesTradingStatusZmq] Error sending trading status %s, error: %+v", tradingStatus, err)
						pub.Close()
						ctx.Term()
						break
					}
				}
			}
		}
	}()
}
