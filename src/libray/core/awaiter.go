package core

import (
	"reflect"
	"time"
)

// 并发异步模型
type Future struct {
	ctx   *ContextTimeout
	await func(ctx Context) (retArr []reflect.Value, reterr error)
}

// 异步等待
func (that *Future) Await() ([]reflect.Value, error) {
	return that.await(Background())
}

// 异步超时等待
func (that *Future) AwaitTimeout(timeout time.Duration) ([]reflect.Value, error) {
	ctx, cancel := WithTimeoutEx(Background(), timeout)
	if ctx == nil || cancel == nil {
		return nil, DeadlineExceeded
	}
	that.ctx = ctx
	defer cancel()
	return that.await(that.ctx)
}

// 设置超时
func (that *Future) SetTimeout(timeout time.Duration) {
	if that.ctx != nil {
		that.ctx.SetTimeout(timeout)
	}
}

// 创建异步模型
func Async(cb interface{}, args ...interface{}) (r1 *Future) {
	kind := reflect.TypeOf(cb).Kind()
	if kind != reflect.Func {
		Logger.Error("cb must be func")
		return
	}

	// 参数处理
	handler := reflect.ValueOf(cb)
	funcType := handler.Type()
	options := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			options[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			options[i] = reflect.ValueOf(v)
		}
	}

	// 执行函数
	taskChan := make(chan struct{})
	var retArr []reflect.Value
	var reterr error
	go func(handler_ reflect.Value, options_ []reflect.Value) {
		defer func() {
			defer close(taskChan)
			if err := recover(); err != nil {
				Logger.Errorf("failed to await %v", err)
			}
		}()
		retArr = handler_.Call(options_)
		if len(retArr) > 0 {
			last := retArr[len(retArr)-1].Interface()
			opterr, ok := last.(error)
			if ok {
				retArr = nil
				reterr = opterr
			}
		}
	}(handler, options)

	// 信号处理
	return &Future{
		await: func(ctx Context) ([]reflect.Value, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-taskChan:
				return retArr, reterr
			}
		},
	}
}
