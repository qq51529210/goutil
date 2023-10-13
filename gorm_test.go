package util

import "testing"

type testInitGORMQuery struct {
	A int64
}

func TestInitGORMQuery(t *testing.T) {
	v := new(testInitGORMQuery)
	InitGORMQuery(nil, v)
}
