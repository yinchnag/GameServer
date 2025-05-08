package network

import (
	"sync"
	"sync/atomic"

	"wgame_server/libray/core"
)

const (
	SESSION_REMOTE_SHIFT = 24         // 节点ID集群便宜
	SESSION_SLOT_SIZE    = 1024       // 存储槽初始大小
	SESSION_ID_MASK      = 0x00ffffff // 节点ID掩码
)

func newWsSessionStorage() *WsSessionStorage {
	inst := &WsSessionStorage{}
	inst.init()
	return inst
}

type WsSessionStorage struct {
	cluster  uint32              // 集群ID
	index    uint32              // 存储索引
	slotSize int                 // 槽位数
	slots    []*WsConnectSession // 存储槽
	lock     sync.RWMutex        // 读写锁
	count    int32               // 客户端数量
}

func (that *WsSessionStorage) init() {
	that.cluster = 1 << SESSION_REMOTE_SHIFT
	that.index = 1
	that.slotSize = SESSION_SLOT_SIZE
	that.slots = make([]*WsConnectSession, that.slotSize)
}

func (that *WsSessionStorage) FindSessionByID(sessionId uint32) *WsConnectSession {
	that.lock.RLock()
	defer that.lock.RUnlock()
	hash := sessionId & uint32(that.slotSize-1)
	find := that.slots[hash]
	if find != nil && find.sessionId == sessionId {
		return find
	}
	return nil
}

func (that *WsSessionStorage) RegisterRecv(client *WsConnectSession) uint32 {
	that.lock.Lock()
	defer that.lock.Unlock()
	for {
		sessionID := uint32(0)
		for k, v := range that.slots {
			sessionID = 1 // 保留0号位
			if v == nil {
				sessionID = that.cluster | uint32(k)
				client.sessionId = sessionID
				that.index = sessionID + 1
				that.slots[k] = client
				atomic.AddInt32(&that.count, 1)
				return sessionID
			}
		}
		if that.slotSize*2 > SESSION_ID_MASK {
			core.Logger.Error("Session storage is full")
			return 0
		}

		slots := make([]*WsConnectSession, that.slotSize*2)
		copy(slots[:that.slotSize], that.slots)
		that.slots = slots
		that.slotSize *= 2
	}
}

func (that *WsSessionStorage) Unregister(client *WsConnectSession) {
	that.lock.Lock()
	defer that.lock.Unlock()
	hash := client.sessionId & uint32(that.slotSize-1)
	find := that.slots[hash]
	if find != nil && find.sessionId == client.sessionId {
		that.slots[hash] = nil
	} else {
		core.Logger.Error("Unregister session error")
	}
}

func (that *WsSessionStorage) GetCount() int32 {
	return atomic.LoadInt32(&that.count)
}

func (that *WsSessionStorage) ForEach(cb func(client *WsConnectSession)) {
	that.lock.RLock()
	defer that.lock.RUnlock()
	for _, client := range that.slots {
		if client != nil {
			cb(client)
		}
	}
}
