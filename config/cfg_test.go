package config

import (
	"fmt"
	"os"
	"testing"
)

type testReadEnv1 struct {
	F1 string
	F2 string `env:"test1"`
	F3 testReadEnv2
}

type testReadEnv2 struct {
	F1 string
	F2 int `env:"test2"`
	F3 testReadEnv3
}

type testReadEnv3 struct {
	F1 string
	F2 float64 `env:"test3"`
}

func Test_ReadEnv(t *testing.T) {
	os.Setenv("test1", "1")
	os.Setenv("test2", "2")
	os.Setenv("test3", "3.1")
	//
	var m testReadEnv1
	// json
	// env
	ReadEnv(&m)
	fmt.Println(m)
}
