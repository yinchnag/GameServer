package database

import (
	"reflect"

	"wgame_server/libray/core"
)

type DbiStorageObj struct {
	kind reflect.Kind    // 数据类型 主要用来区分 struct map  (slice,array)
	tag  []string        // 标签列表
	val  []reflect.Value // 值列表
}

func (that *DbiStorageObj) Encode(inst any) {
	if len(that.tag) != 0 { // 已经映射过一次了
		return
	}
	defer func() {
		if err := recover(); err != nil {
			core.Logger.Errorln("Encode 拦截到错误:", err)
		}
	}()
	dstType := reflect.TypeOf(inst).Elem() // 获得类型
	that.kind = dstType.Kind()
	rawData := reflect.ValueOf(inst).Elem() // 获得值
	for i := 0; i < dstType.NumField(); i++ {
		tag := dstType.Field(i).Tag.Get("json")
		name := dstType.Field(i).Name
		if tag == "" {
			continue
		}
		ignore := dstType.Field(i).Tag.Get("ignore") // s是否忽略
		if ignore == "1" {
			continue
		}
		if idx := core.FindSlice(that.tag, func(val string, idx int) bool {
			return val == tag
		}); idx != -1 {
			core.Logger.Warn("重复映射", name)
			continue
		}
		that.tag = append(that.tag, tag)
		that.val = append(that.val, rawData.Field(i))
	}
}

func (that *DbiStorageObj) FromJson(jsonData string) {
}

func (that *DbiStorageObj) ToJson() {
}
