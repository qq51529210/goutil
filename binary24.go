package util

// PutBigUint24 写入 3 字节
func PutBigUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// BigUint24 读取 3 字节
func BigUint24(b []byte) uint32 {
	v := uint32(b[0]) << 16
	v |= uint32(b[1]) << 8
	v |= uint32(b[2])
	return v
}

// PutLittleUint24 写入 3 字节
func PutLittleUint24(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

// LittleUint24 读取 3 字节
func LittleUint24(b []byte) uint32 {
	v := uint32(b[0])
	v |= uint32(b[1]) << 8
	v |= uint32(b[2]) << 16
	return v
}
