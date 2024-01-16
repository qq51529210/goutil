package device

import (
	"context"
	"fmt"
	"testing"
)

const (
	host     = "192.168.31.66"
	username = "admin"
	password = "hwonvif66"
)

func Test_GetSystemDateAndTime(t *testing.T) {
	d, err := NewDevice(context.Background(), host, username, password)
	if err != nil {
		t.Fatal(err)
	}
	m, err := d.GetSystemDateAndTime(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m.UTC())
	fmt.Println(m.Local())
}

func Test_GetDeviceInformation(t *testing.T) {
	// d, err := NewDevice(context.Background(), host, username, password)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(m)
}
