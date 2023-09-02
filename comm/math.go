package comm

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Abs[T Number](a, b T) T {
	if a < b {
		return b - a
	} else {
		return a - b
	}
}
