// JoysGames copyrights this specification. No part of this specification may be
// reproduced in any form or means, without the prior written consent of JoysGames.
//
// This specification is preliminary and is subject to change at any time without notice.
// JoysGames assumes no responsibility for any errors contained herein.
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
// @package JGServer
// @copyright joysgames.cn All rights reserved.
// @version v1.0

package core

import (
	"time"
)

// 过期上下文
// 不保证修改后子节点顺利结束
func WithTimeoutEx(parent Context, timeout time.Duration) (*ContextTimeout, CancelFunc) {
	return WithDeadlineEx(parent, time.Now().Add(timeout))
}

// 过期上下文
// 不保证修改后子节点顺利结束
func WithDeadlineEx(parent Context, date time.Time) (*ContextTimeout, CancelFunc) {
	if parent == nil {
		return nil, nil
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(date) {
		return nil, nil
	}
	ctx := &ContextTimeout{
		timerCtx{
			cancelCtx: newCancelCtx(parent),
			deadline:  date,
		},
	}
	propagateCancel(parent, ctx)
	dur := time.Until(date)
	if dur <= 0 {
		ctx.cancel(true, DeadlineExceeded) // deadline has already passed
		return ctx, func() { ctx.cancel(false, ErrCanceled) }
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.err == nil {
		ctx.timer = time.AfterFunc(dur, func() {
			ctx.cancel(true, DeadlineExceeded)
		})
	}
	return ctx, func() { ctx.cancel(true, ErrCanceled) }
}

// 扩展倒计时上下文
type ContextTimeout struct {
	timerCtx
}

// 充值倒计时
func (that *ContextTimeout) SetTimeout(timeout time.Duration) {
	that.deadline = time.Now().Add(timeout)
	dur := time.Until(that.deadline)
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.timer != nil {
		that.timer.Stop()
		that.timer = nil
	}
	that.timer = time.AfterFunc(dur, func() {
		that.cancel(true, DeadlineExceeded)
	})
}
