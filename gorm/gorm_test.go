package gorm

import (
	"fmt"
	"goutil/log"
	"testing"
)

type InitMysqlJSONSetT struct {
	ID int64
	JF []byte `gorm:"type:json"`
	DL string `gorm:"type:varchar(12)"`
}

type InitMysqlJSONSetTJF1 struct {
	// 测试不可导出
	f1 string
	// 正常测试
	F1 int
	// 测试忽略零值
	F2 string `json:",omitempty"`
	// 测试忽略
	F3 string `json:"-"`
	// 测试名称
	F4 string `json:"ff"`
	// 测试嵌入
	InitMysqlJSONSetJF2
	// 测试结构
	F5 InitMysqlJSONSetJF2
	// 空指针，不处理
	F6 *InitMysqlJSONSetJF2
}

type InitMysqlJSONSetJF2 struct {
	// 嵌入同名，直接覆盖
	F1 string
	F2 int64 `json:",omitempty"`
}

func Test_InitMysqlJSONSet(t *testing.T) {
	db, err := Open("mysql://root:mysql1234abcd@tcp(192.168.58.155:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	//
	m := &InitMysqlJSONSetT{}
	if err := db.AutoMigrate(m); err != nil {
		t.Fatal(err)
	}
	db.Config.Logger = NewLog(log.DefaultLogger)
	//
	var res PageResult[*InitMysqlJSONSetT]
	if err := Page(db, &PageQuery{}, &res, "ID", "DL"); err != nil {
		t.Fatal(err)
	}
	//
	// {
	// 	err = db.AutoMigrate(m)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	m.ID = 4
	// 	m.JF = []byte(`{"a1":1,"F5":{},"F6":{"a1":22}}`)
	// 	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error; err != nil {
	// 		t.Fatal(err)
	// 	}
	// }
	// {
	// 	d1 := &InitMysqlJSONSetTJF1{}
	// 	// d1.f1 = "f1"
	// 	d1.F1 = 1
	// 	// d1.F3 = "f3"
	// 	// d1.F4 = "f4"
	// 	d1.InitMysqlJSONSetJF2.F1 = "conv"
	// 	d1.F5.F1 = "f5.f1"
	// 	d1.F5.F2 = 123
	// 	if err := db.Model(m).UpdateColumn("JF", InitMysqlJSONSet(d1, "JF")).Error; err != nil {
	// 		t.Fatal(err)
	// 	}
	// }
	fmt.Println()
}
