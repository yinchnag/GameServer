package actor

import (
	"sync"

	"wgame_server/libray/core"
)

func newActorStorage(cluster uint32) *ActorStorage {
	inst := &ActorStorage{}
	inst.init(cluster)
	return inst
}

type ActorStorage struct {
	cluster  uint32          // 集群ID
	index    uint32          // 存储索引
	slotSize int             // 槽位数
	slots    []*ActorContext // 槽位
	lock     sync.RWMutex    // 读写锁
}

func (that *ActorStorage) init(cluster uint32) {
	that.cluster = (cluster & 0xff) << ACTORID_REMOTE_SHIFT
	that.index = 1
	that.slotSize = ACTORID_SLOT_SIZE
	that.slots = make([]*ActorContext, that.slotSize)
}

func (that *ActorStorage) findActorByID(actorID uint32) *ActorContext {
	that.lock.RLock()
	defer that.lock.RUnlock()
	for _, v := range that.slots {
		if v != nil && v.ActorID() == actorID {
			return v
		}
	}
	return nil
}

func (that *ActorStorage) findActorByAlias(alias string) *ActorContext {
	that.lock.RLock()
	defer that.lock.RUnlock()
	for _, v := range that.slots {
		if v != nil && v.Alias() == alias {
			return v
		}
	}
	return nil
}

func (that *ActorStorage) toActorID(alias string) uint32 {
	that.lock.RLock()
	defer that.lock.RUnlock()
	for _, v := range that.slots {
		if v != nil && v.Alias() == alias {
			return v.ActorID()
		}
	}
	return 0
}

func (that *ActorStorage) register(ctx *ActorContext) uint32 {
	that.lock.Lock()
	defer that.lock.Unlock()
	for {
		actorID := uint32(0)
		for k, v := range that.slots {
			if actorID > ACTORID_MASK {
				actorID = 1 // 保留0号
			}
			if v == nil {
				ctx.runner.actorID = actorID
				that.index = actorID + 1
				that.slots[k] = ctx
				return actorID
			}
		}
		if that.slotSize*2 > ACTORID_MASK {
			core.Logger.Error("actor storage is full")
			return 0
		}

		slots := make([]*ActorContext, that.slotSize*2)
		copy(slots[:that.slotSize], that.slots)
		that.slots = slots
		that.slotSize *= 2
	}
}

func (that *ActorStorage) ForEach(cb func(*ActorContext)) {
	that.lock.RLock()
	defer that.lock.RUnlock()
	for _, ctx := range that.slots {
		if ctx != nil {
			cb(ctx)
		}
	}
}
