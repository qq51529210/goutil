package discovery

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	d, err := NewDiscoverer("en0", MulticastAddr, func(addr string, err error) {
		if err != nil {
			fmt.Println(err)
		}
		if addr != "" {
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(u.Host)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for i := 0; i < 2; i++ {
		if err := d.Discover(time.Second * 5); err != nil {
			t.Fatal(err)
		}
	}
}
