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
	password = "aaa123123"
)

func Test_GetSystemDateAndTime(t *testing.T) {
	d := NewDevice(host, username, password, 0)
	m, err := d.GetSystemDateAndTime(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m.UTC())
	fmt.Println(m.Local())
}

func Test_GetGetCapabilities(t *testing.T) {
	d := NewDevice(host, username, password, 0)
	m, err := d.GetCapabilities(context.Background(), owvd.CapabilityCategoryAll)
	if err != nil {
		t.Fatal(err)
	}
	if m.Analytics != nil {
		fmt.Println("Analytics", m.Analytics.XAddr)
	}
	if m.Device != nil {
		fmt.Println("Device", m.Device.XAddr)
	}
	if m.Events != nil {
		fmt.Println("Events", m.Events.XAddr)
	}
	if m.Imaging != nil {
		fmt.Println("Imaging", m.Imaging.XAddr)
	}
	if m.Media != nil {
		fmt.Println("Media", m.Media.XAddr)
	}
	if m.PTZ != nil {
		fmt.Println("PTZ", m.PTZ.XAddr)
	}
	if m.Extension != nil {
		if m.Extension.DeviceIO != nil {
			fmt.Println("DeviceIO", m.Extension.DeviceIO.XAddr)
		}
		if m.Extension.Display != nil {
			fmt.Println("Display", m.Extension.Display.XAddr)
		}
		if m.Extension.Recording != nil {
			fmt.Println("Recording", m.Extension.Recording.XAddr)
		}
		if m.Extension.Search != nil {
			fmt.Println("Search", m.Extension.Search.XAddr)
		}
		if m.Extension.Replay != nil {
			fmt.Println("Replay", m.Extension.Replay.XAddr)
		}
		if m.Extension.Receiver != nil {
			fmt.Println("Receiver", m.Extension.Receiver.XAddr)
		}
		if m.Extension.AnalyticsDevice != nil {
			fmt.Println("AnalyticsDevice", m.Extension.AnalyticsDevice.XAddr)
		}
	}
}

func Test_GetDeviceInformation(t *testing.T) {
	d := NewDevice(host, username, password, 0)
	m, err := d.GetDeviceInformation(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}

func Test_GetProfiles(t *testing.T) {
	d, err := NewDeviceWithCapabilities(context.Background(), host, username, password, 0)
	if err != nil {
		t.Fatal(err)
	}
	m, err := d.GetProfiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}

func Test_GetStreamURL(t *testing.T) {
	d, err := NewDeviceWithCapabilities(context.Background(), host, username, password, 0)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := d.GetProfiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range ps {
		m, err := d.GetStreamURL(context.Background(), p.Token)
		if err != nil {
			t.Fatal(err)
		}
		m.URL = d.AuthStreamURL(m.URL)
		fmt.Println(*m)
	}
}
