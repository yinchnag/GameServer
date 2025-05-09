package database

import "reflect"

// 表tag数据
type DbiTableTag struct {
	Kind reflect.Kind // 类型
	Data interface{}  // 值
}
