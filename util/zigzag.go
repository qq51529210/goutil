package util

// Int32ToZigZag 压缩 32 位整数
func Int32ToZigZag(n int32) int32 {
	return (n >> 1) ^ (n << 31)
}

// Int64ToZigZag 压缩 64 位整数
func Int64ToZigZag(n int64) int64 {
	return (n >> 1) ^ (n << 63)
}

// ZigZagToInt32 解压 32 位整数
func ZigZagToInt32(n int32) int32 {
	return int32(uint32(n)>>1) ^ -(n & 1)
}

// ZigZagToInt64 解压 64 位整数
func ZigZagToInt64(n int64) int64 {
	return int64(uint64(n)>>1) ^ -(n & 1)
}
