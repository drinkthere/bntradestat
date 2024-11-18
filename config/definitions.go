package config

type (
	TradingStatus string
)

const (
	TradingStatusOK   = TradingStatus("OK")
	TradingStatusStop = TradingStatus("STOP")
)
