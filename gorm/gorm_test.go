package gorm

import (
	"encoding/json"
	"testing"
)

type InitMysqlJSONSetT struct {
	ID int64
	JF []byte `gorm:"type:json"`
}

type InitMysqlJSONSet1 struct {
	S2 InitMysqlJSONSet2 `json:",omitempty"`
	InitMysqlJSONSet2
	F1 int    `json:""`
	F2 string `json:",omitempty"`
	F3 string `json:"-"`
	F4 string `json:"F3,omitempty"`
}

type InitMysqlJSONSet2 struct {
	F1 string `json:",omitempty"`
	F2 float64
}

func Test_InitMysqlJSONSet(t *testing.T) {
	db, err := Open("mysql://root:mysql1234abcd@tcp(192.168.58.155:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	//
	m := &InitMysqlJSONSetT{}
	err = db.AutoMigrate(m)
	if err != nil {
		t.Fatal(err)
	}
	m1 := &InitMysqlJSONSet1{}
	m1.S2.F1 = "s2.f1"
	m1.S2.F2 = 2.2
	m1.InitMysqlJSONSet2.F1 = "f1"
	m1.InitMysqlJSONSet2.F2 = 2
	m1.F1 = 1
	m1.F2 = "s1.f1"
	m.JF, _ = json.Marshal(m1)
	//
	if err := db.Create(m).Error; err != nil {
		t.Fatal(err)
	}
	//
	m2 := &InitMysqlJSONSet1{}
	// m2.S2.F1 = "s2.f12"
	// m2.S2.F2 = 22.2
	// m2.InitMysqlJSONSet2.F1 = "f12"
	// m2.F4 = "s1.4223"
	m2.F3 = "12312"
	m2.InitMysqlJSONSet2.F2 = 22
	// m2.F1 = 12
	m2.F2 = "s1.f123123"
	//
	m.ID = 1
	if err := db.Model(m).UpdateColumn("JF", InitMysqlJSONSet(m2, "`JF`")).Error; err != nil {
		t.Fatal(err)
	}
}
