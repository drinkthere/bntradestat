package container

import (
	"bntradestat/config"
	"sort"
	"sync"
	"time"
)

type TradingStat struct {
	Volume float64
	Times  int
}

type VolumeWrapper struct {
	Volume      float64
	TradeTimeMs int64 //更新时间（毫秒）
}

type TradingComposite struct {
	tradingStatus     config.TradingStatus
	volumeList        []VolumeWrapper
	rwLock            *sync.RWMutex
	lastAggTradeIDMap map[string]int64
}

func NewTradingComposite(globalConfig *config.Config) *TradingComposite {
	composite := TradingComposite{
		tradingStatus:     config.TradingStatusOK,
		volumeList:        make([]VolumeWrapper, 0),
		rwLock:            new(sync.RWMutex),
		lastAggTradeIDMap: make(map[string]int64),
	}
	for _, instID := range globalConfig.InstIDs {
		composite.lastAggTradeIDMap[instID] = 0
	}

	return &composite
}

func (tc *TradingComposite) UpdateTradingInfo(instID string, tradeID int64, volume VolumeWrapper, rollingPeriod int) {
	tc.rwLock.Lock()
	defer tc.rwLock.Unlock()
	// 防止重复更新
	if tradeID <= tc.lastAggTradeIDMap[instID] {
		return
	}
	// 更新最新的trade id
	tc.lastAggTradeIDMap[instID] = tradeID

	// 添加新数据
	tc.volumeList = append(tc.volumeList, volume)

	// 删除过期数据
	now := time.Now().UnixNano() / 1e6

	// 删除 volumeList 中过期的数据
	// 使用二分查找找到第一个在时间段中的数据，既tradeTimeMs >= now-int64(rollingPeriod*1000) 的元素
	// 然后将前面的元素都视为过期元素，并删除
	idx := sort.Search(len(tc.volumeList), func(i int) bool {
		return tc.volumeList[i].TradeTimeMs >= now-int64(rollingPeriod*1000)
	})

	tc.volumeList = tc.volumeList[idx:]
}

func (tc *TradingComposite) StatTradingInfo() TradingStat {
	tc.rwLock.RLock()
	defer tc.rwLock.RUnlock()

	var totalQuantity float64
	for _, v := range tc.volumeList {
		totalQuantity += v.Volume
	}

	return TradingStat{
		Volume: totalQuantity,
		Times:  len(tc.volumeList),
	}
}

func (tc *TradingComposite) GetTradingStatus() config.TradingStatus {
	tc.rwLock.RLock()
	defer tc.rwLock.RUnlock()
	return tc.tradingStatus
}

func (tc *TradingComposite) UpdateTradingStatus(status config.TradingStatus) {
	tc.rwLock.Lock()
	defer tc.rwLock.Unlock()
	tc.tradingStatus = status
}
