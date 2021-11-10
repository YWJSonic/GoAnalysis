package structtest

import (
	wwe "codeanalysis/load/project/goloader"
	"sync"
)

type G16Options struct {
	writeFlow chan *wwe.ANAT1 // 提供使用者發布訊息的通道
	ww        chan<- string
	// w2w             chan string
	writeFlowBackup <-chan *wwe.ANAT1 // 暫存publish失敗的訊息，若重新連線成功再次發送
	writeEnd        chan struct{}     // 用於停止 publish() 生命週期
	shutdownChan    chan struct{}     // graceful shutdown 信號

	subscriberMap map[string]func() error // map[RoomGuid]func() error 記住所有註冊方法
	mu            sync.Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	Antes            []wwe.GoFileNode `json:"antes"`            // 底注
	MinCoins         []wwe.GoFileNode `json:"minCoins"`         // 最低金額
	Rounds           []wwe.GoFileNode `json:"rounds"`           // 局數
	DisabledPatterns []wwe.GoFileNode `json:"disabledPatterns"` // 缺少的獎項
	BirdCatchCounts  []wwe.GoFileNode `json:"birdCatchCounts"`  // 抓鳥牌型
}
