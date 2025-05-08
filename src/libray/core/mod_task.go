package core

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"

	"github.com/timandy/routine"
)

// 协程buff池
var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

// 获取协程id(调用10000次50ms)
func GetGoroutineID() int64 {
	goroutineId := routine.Goid()
	if goroutineId == 0 {
		bp := littleBuf.Get().(*[]byte)
		defer littleBuf.Put(bp)
		buff := *bp
		runtime.Stack(buff, false)
		buff = bytes.TrimPrefix(buff, []byte("goroutine "))
		buff = buff[:bytes.IndexByte(buff, ' ')]
		goroutineId, _ = strconv.ParseInt(string(buff), 10, 64)
	}
	return goroutineId
}
