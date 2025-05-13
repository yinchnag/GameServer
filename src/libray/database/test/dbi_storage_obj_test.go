package test

import (
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
)

type StorageItem struct {
	ItemId    int64            `json:"item"`
	Count     int64            `json:"count"`
	Timestamp int64            `json:"ignore"`
	Slice     []int64          `json:"slice"`
	Map       map[string]int64 `json:"Map"`
	Itemd     *Item            `json:"itemd"`
	// database.DbiStorageObj
}

func TestStorage(t *testing.T) {
	s := &StorageItem{
		ItemId: 1,
		Count:  1,
		Slice:  []int64{1, 2, 3, 4, 5},
		Map:    map[string]int64{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
		Itemd:  &Item{ItemId: 1, Count: 1},
	}
	str, _ := sonic.Marshal(s)
	s1 := &StorageItem{}
	sonic.Unmarshal(str, s1)
	fmt.Printf(string(str))
}

func BenchmarkSonicMarshalAndUnmarshal(b *testing.B) {
	s := &StorageItem{
		ItemId: 1,
		Count:  1,
		Slice:  []int64{1, 2, 3, 4, 5},
		Map:    map[string]int64{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
		Itemd:  &Item{ItemId: 1, Count: 1},
	}
	s1 := &StorageItem{}
	for i := 0; i < b.N; i++ {
		b, _ := sonic.Marshal(s)

		sonic.Unmarshal(b, s1)
	}
}
