package actor

const (
	ACTORID_OVERLOAD     = 1024
	ACTORID_SUSPEND      = 10
	ACTORID_REMOTE_SHIFT = 24
	ACTORID_SLOT_SIZE    = 4
	ACTORID_MASK         = 0xffffff
)

// 运行权重
var RUNNER_WEIGHT = [32]int{
	-1, -1, -1, -1, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3,
}

type actorWorker struct {
	index   int           // 索引
	context *ActorContext // 当前上下文
	cursor  uint32        // 游标
}

type actorSuspend struct {
	ctx *ActorContext
}
