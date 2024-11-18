package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Source struct {
	IP   string // local IP to connect to OKx websocket
	Colo bool   // is co-location with Okx
}

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// 币安配置
	BinanceAPIKey    string
	BinanceSecretKey string

	// 获取trade消息的source ip以及是否使用内网
	Sources []Source

	// 订阅的交易对
	InstIDs []string

	// 价格最小精度
	MinAccuracy float64

	// 统计trade data的有效时间范围，以s记
	RollingPeriodSeconds int
	// 最小交易次数
	MinTradingTimes int
	// 最小交易
	MinTradingVolume float64

	// 发送交易状态的ZMQ
	TradingStatusZMQIPC string
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
