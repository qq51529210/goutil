package uid

import (
	"crypto/md5"
	"encoding/binary"
	"strconv"
)

// SetSeed 对 seed 进行 md5 得到 16 位整数
// 然后设置 uuid node (长度不够前面补 0 ) 和 snowfalke mechine id
func SetSeed(seed string) {
	// 随机数
	h := md5.New()
	h.Write([]byte(seed))
	h.Sum(nil)
	id := binary.BigEndian.Uint16(h.Sum(nil))
	s := strconv.Itoa(int(id))
	for len(s) < 6 {
		s = "0" + s
	}
	SetUUIDNode(s)
	SetSnowflakeMechineID(id % 4095)
}
