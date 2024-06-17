package strings

// KeyValue 解析 k=v
func KeyValue(line string) (string, string) {
	return KeyValueWith(line, CharEqual)
}

// KeyValueWith 解析 k(split)v
func KeyValueWith(line string, split byte) (string, string) {
	// k = v
	k, v := Split(line, split)
	// k 去空白
	k = TrimByte(k, CharSpace, CharSpace)
	// v 去空白
	v = TrimByte(v, CharSpace, CharSpace)
	//
	return k, v
}
