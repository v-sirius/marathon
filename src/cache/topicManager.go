package cache

import (
	"sync"
)

// 位置相关topic
// 位置topic的订阅关键字是用户唯一代码
type PositionTopic struct {
	Topic2Exist map[string]bool
	rwLock      sync.RWMutex
}

var G_PositionTopic = new(PositionTopic)

// 初始化位置topic
func (pPosTopic *PositionTopic) InitPositionTopic() {
	pPosTopic.Topic2Exist = make(map[string]bool, 0)
}

// 添加用户唯一代码
func (pPosTopic *PositionTopic) AddPositionTopic(ucode string) {
	pPosTopic.rwLock.Lock()
	defer pPosTopic.rwLock.Unlock()

	pPosTopic.Topic2Exist[ucode] = true
}

// 取消订阅
func (pPosTopic *PositionTopic) DelPositioonTopic(ucode string) {
	pPosTopic.rwLock.Lock()
	defer pPosTopic.rwLock.Unlock()

	_, exist := pPosTopic.Topic2Exist[ucode]
	if exist {
		pPosTopic.Topic2Exist[ucode] = false
	}
}

// 检查是否已订阅用户唯一标识为ucode的位置信息
func (pPosTopic PositionTopic) IsSubscribed(ucode string) bool {
	pPosTopic.rwLock.RLock()
	defer pPosTopic.rwLock.RUnlock()

	_, exist := pPosTopic.Topic2Exist[ucode]
	if exist && pPosTopic.Topic2Exist[ucode] == true {
		return true
	}

	return false
}
