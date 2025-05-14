package core

// 三目运算的函数
func Ternary[T comparable](check bool, b T, c T) T {
	if check {
		return b
	}
	return c
}

// 三目运算的函数
func TernaryF[T comparable](check bool, b func() T, c T) (r1 T) {
	if check {
		r1 = b()
	} else {
		r1 = c
	}
	return r1
}

// 三目运算的函数
func TernaryFF[T comparable](check bool, b func() T, c func() T) (r1 T) {
	if check {
		r1 = b()
	} else {
		r1 = c()
	}
	return r1
}
