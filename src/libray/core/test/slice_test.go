package test

import (
	"testing"
	"wgame_server/libray/core"
)

func TestSliceFind(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	index := core.FindSlice(arr, func(value int, index int) bool {
		return value == 2
	})
	if index == -1 {
		t.Fail()
	}
}

func BenchmarkSliceFind(b *testing.B) {
	arr := []int{}
	for i := 0; i < 200; i++ {
		arr = append(arr, i)
	}
	for i := 0; i < b.N; i++ {
		// index := core.FindSlice(arr, func(value int, index int) bool {
		// 	return value == 123
		// })
		// _ = index
		for _, v := range arr {
			if v == 123 {
				break
			}
		}
	}
}

func TestSliceRemove(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6}
	arr = core.SliceRemoveByIndex(arr, 2)
	if len(arr) != 5 {
		t.Fail()
	}
}
