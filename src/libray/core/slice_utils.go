package core

import "reflect"

// 查询效率比直接调用for循环遍历低了4倍 在i7-9750H 2.6GHz配置下
func FindSlice[T Any](arr []T, cb func(value T, index int) bool) int {
	if len(arr) == 0 || cb == nil {
		return -1
	}
	for index, value := range arr {
		if cb(value, index) {
			return index
		}
	}
	return -1
}

// 根据下标删除切片元素
func SliceRemoveByIndex[T Any](arr []T, index int) []T {
	size := len(arr)
	if index < 0 || index >= size {
		return nil
	}
	arr = append(arr[:index], arr[index+1:]...)
	return arr
}

func SliceRemoveByVal(arr []reflect.Value, value reflect.Value) []reflect.Value {
	if len(arr) == 0 {
		return arr
	}
	index := -1
	for _, ele := range arr {
		index++
		if ele == value {
			break
		}
	}
	arr = append(arr[:index], arr[index+1:]...)
	return arr
}
