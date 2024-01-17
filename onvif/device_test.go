package onvif

import (
	"context"
	"fmt"
	"testing"

	owvd "goutil/onvif/wsdl/ver10/device"
)

const (
	// host     = "192.168.31.66"
	// username = "admin"
	// password = "hwonvif66"
	host     = "192.168.31.3"
	username = "ovfs"
	password = "aa123123"
	// username = "admin"
	// password = "A@12345678"
)

func Test_GetSystemDateAndTime(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := owvd.GetSystemDateAndTime(context.Background(), d.URL)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m.UTC())
	fmt.Println(m.Local())
}

func Test_GetGetCapabilities(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := owvd.GetCapabilities(context.Background(), d.URL, d.Security, owvd.CapabilityCategoryAll)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}

func Test_GetDeviceInformation(t *testing.T) {
	d := NewDevice(context.Background(), host, username, password)
	m, err := owvd.GetDeviceInformation(context.Background(), d.URL, d.Security)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}
