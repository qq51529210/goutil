package reflect

import (
	"fmt"
	"testing"
)

type structCopySrc struct {
	// int
	StructCopys
	// *StructCopys
	// 不导出结构
	structCopys
	// *structCopys
	// string
	// src 有值也忽略
	Src1 string `cp:"-"`
	// src 指定字段
	Src2 string `cp:"Dst2"`
	// src 忽略零值
	Src3 int `cp:"omitempty"`
	// src 值，dst 指针
	Src4 float64
	// src 指针，dst 值
	Src5 *int
	// src 可导出，dst 不导出
	Src6 string `cp:"dst6"`
	// src 不导出，dst 可导出
	src7 string `cp:"Dst7"`
	// nil 指针覆盖
	Dst8 *int
}

type structCopyDst struct {
	StructCopys
	// *StructCopys
	structCopys
	// *structCopys
	//
	Src1 string
	Src2 string
	Dst2 string
	Src3 *int
	Src4 *float64
	Src5 int
	dst6 string
	Dst7 string
	Dst8 *int
}

type structCopys struct {
	S1 *int
	S2 int
}

type StructCopys struct {
	SS1 *int
	SS2 int
}

func Test_StructCopy(t *testing.T) {
	var src structCopySrc
	//
	src.StructCopys.SS1 = new(int)
	*src.StructCopys.SS1 = 11
	src.StructCopys.SS2 = 1122
	//
	src.structCopys.S1 = new(int)
	*src.structCopys.S1 = 22
	src.structCopys.S2 = 2222
	//
	src.Src1 = "src1"
	src.Src2 = "src2"
	src.Src3 = 3
	src.Src4 = 0.4
	src.Src5 = new(int)
	*src.Src5 = 5
	src.Src6 = "src6"
	//
	var dst structCopyDst
	dst.Dst8 = new(int)
	StructCopy(&dst, src, "cp")
	//
	fmt.Println(dst)
}
