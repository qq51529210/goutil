package util

func MaxIn[T int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](arr []T) T {
	var n T
	if len(arr) < 1 {
		return n
	}
	// 查找
	n = arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > n {
			n = arr[i]
		}
	}
	// 返回
	return n
}

func MinIn[T int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](arr []T) T {
	var n T
	if len(arr) < 1 {
		return n
	}
	// 查找
	n = arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] < n {
			n = arr[i]
		}
	}
	// 返回
	return n
}
