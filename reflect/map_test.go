package reflect

import (
	"encoding/json"
	"fmt"
	"testing"
)

type StructFromMap1 struct {
	F11 string
	F12 *string
	F13 *float64 `fm:"ff"`
	F14 int      `fm:"-"`
	*StructFromMap2
	StructFromMap3
	F15 *StructFromMap2 `fm:"-"`
	F16 StructFromMap2  `fm:"f6"`
	F17 *StructFromMap2
}

type StructFromMap2 struct {
	F21 string
	F22 *string
}

type StructFromMap3 struct {
	F31 string
	F32 *string
}

func Test_StructFromMap(t *testing.T) {
	m := make(map[string]any)
	m["F11"] = "a1"
	m["F12"] = "a2"
	m["F13"] = "a3"
	m["ff"] = 1.9
	m["F14"] = 3
	m["F21"] = "21a"
	m["F22"] = "22a"
	m["F22"] = "22a"
	m["F22"] = "22a"
	m["F15"] = map[string]any{
		"F21": "1",
		"F22": "1",
	}
	m["f6"] = map[string]any{
		"F21": "1",
		"F22": "1",
	}
	m["F17"] = map[string]any{
		"F21": "21",
		"F22": "21",
	}
	s := new(StructFromMap1)
	StructFromMap(s, "fm", m)
	fmt.Println(s)
}

type structToMapS1 struct {
	string
	S1F1 string  `tm:"s1f1"`
	S1F2 *string `tm:"omitempty"`
	S1F3 string  `tm:"omitempty"`
	S1F4 *string `tm:"-"`
	S1F5 any
	s1F6 string
}

type structToMapS2 struct {
	S2F1 int
	S2F2 *string
}

type StructToMapS3 struct {
	S3F1 int
	S3F2 *string
}

type structToMapS struct {
	string
	*int
	A *int
	B int
	structToMapS1
	*structToMapS2
	*StructToMapS3
	SF1 *StructToMapS3 `tm:"omitempty"`
	SF2 structToMapS2
}

func Test_StructToMap(t *testing.T) {
	var s structToMapS
	//
	s.string = "0"
	s.S1F1 = `tm:"s1f1"`
	s.S1F4 = new(string)
	*s.S1F4 = `tm:"-"`
	s.S1F5 = 123.123
	s.s1F6 = "6"
	//
	s.structToMapS2 = new(structToMapS2)
	s.structToMapS2.S2F1 = 111
	s.structToMapS2.S2F2 = new(string)
	*s.structToMapS2.S2F2 = "s."
	//
	s.SF2.S2F1 = 13333
	s.SF2.S2F2 = new(string)
	*s.SF2.S2F2 = "sf2."
	//
	d, _ := json.Marshal(s)
	fmt.Println(string(d))
	m := StructToMap(s, "tm")
	fmt.Println(m)
}
