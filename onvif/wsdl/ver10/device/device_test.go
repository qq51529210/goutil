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
	// username = "onvif"
	// password = "dhonvif3"
)

func Test_GetSystemDateAndTime(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := d.GetSystemDateAndTime(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m.UTC())
	fmt.Println(m.Local())
}

func Test_GetGetCapabilities(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := d.GetCapabilities(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}

func Test_GetDeviceInformation(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := d.GetDeviceInformation(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}
