package core

type Map struct {
	val map[int]int
}

func NewMap() *Map {
	return &Map{
		val: make(map[int]int),
	}
}
