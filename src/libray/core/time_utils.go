package core

import "time"

const DATE_FORMAT1 = "2006-01-02 15:04:05"

// 获得服务器时间
func ServerTime() time.Time {
	return time.Now()
}

// 将时间戳转换成时间
func TimestampToTime(date int64) time.Time {
	return time.Unix(date, 0)
}
