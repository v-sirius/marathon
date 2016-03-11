package cache

import (
	"time"
)

// 用户行为
type Behaviour struct {
	Id        int       // 行为唯一标识
	Operation int       // 操作类型
	OpTime    time.Time // 用户操作时间
	UserId    int       // 用户id
}

// 所有用户行为
type BehaviourSet struct {
	AllBehaviour []Behaviour
	//bFlg         chan bool
}

// 初始化

func InitBehaviourSet() *BehaviourSet {
	pRet := new(BehaviourSet)
	pRet = &BehaviourSet{make([]Behaviour, 0)}

	return pRet
}
