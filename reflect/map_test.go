package reflect

import (
	"encoding/json"
	"fmt"
	"testing"
)

type stFromMap struct {
	// 测试不可导出
	s1 string
	// 测试名称
	S1 string `map:"ss"`
	// 测试零值
	S2 int64
	// 测试数据类型
	S3 int64
	// 测试忽略
	S4 string `map:"-"`
	// 测试空指针自动 new
	S5 *string
	// 测试 nil 不影响指针
	S6 *string
	// 测试嵌入
	*SstFromMap1
	// 测试结构字段
	S7 *SstFromMap1
}

type SstFromMap1 struct {
	SstFromMap2
	// 同名嵌入，类型不同，不会影响
	S1 int64 `map:"ss"`
	// 同名嵌入，类型相同，会被赋值
	S3 string
}

type SstFromMap2 struct {
	S4 int64
	S5 string
	S3 string
}

func Test_StructFromMap(t *testing.T) {
	data := make(map[string]any)
	data["s1"] = "ss"
	data["ss"] = "s1"
	data["S3"] = "s3"
	data["S4"] = "s4"
	data["S5"] = "s5"
	data["S6"] = nil
	data["S7"] = map[string]any{
		"S1": 123,
		"S3": "S333",
		"S4": 412,
		"S5": "s55",
	}
	var s stFromMap
	s.S6 = new(string)
	StructFromMap(&s, data)
	//
	d, _ := json.MarshalIndent(&s, "", " ")
	fmt.Println(string(d))
}

type stToMap struct {
	// 测试不可导出
	s1 string
	// 测试名称
	S1 string `map:"ss"`
	// 测试零值
	S2 int64
	// 测试忽略零值
	S3 int64 `map:",omitempty"`
	// 测试忽略
	S4 string `map:"-"`
	// 不是空指针就不算零值
	S5 *string
	S6 *string `map:",omitempty"`
	// 测试嵌入
	SstToMap1
	// *SstToMap
	// 测试嵌入其他
	int
	// 测试结构字段
	S7 SstToMap1
	// 测试空指针 null
	S8 *SstToMap1
}

type SstToMap1 struct {
	// 测试嵌入名称覆盖，不会
	S1 string `map:",omitempty"`
	// 测试嵌入覆盖
	S2 string `map:",omitempty"`
	// 三层
	S3 SstToMap2
}

type SstToMap2 struct {
	S1 string `map:",omitempty"`
	S2 string
}

func Test_StructToMap(t *testing.T) {
	var s stToMap
	s.s1 = "s"
	s.S1 = "S1"
	s.S2 = 0
	s.S3 = 0
	s.S4 = "s4"
	s.S5 = new(string)
	s.S6 = nil
	// s.SstToMap.S1 = "s.s1"
	// s.SstToMap.S2 = "s.s2"
	s.int = 123
	s.S7.S1 = "s7.s1"
	s.S7.S3.S1 = "s7.s3.s1"
	// s.S7.S2 = "s7.s2"
	data := StructToMap(&s)
	d, _ := json.MarshalIndent(data, "", " ")
	fmt.Println(string(d))
}
