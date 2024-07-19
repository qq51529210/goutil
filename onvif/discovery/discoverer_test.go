package discovery

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	if err := Discover(context.Background(), "en0", MulticastAddr, func(addr string, err error) {
		if err != nil {
			fmt.Println(err)
		}
		if addr != "" {
			for _, a := range strings.Fields(addr) {
				u, err := url.Parse(a)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(a, u.Host)
			}
		}
	}, time.Second*5); err != nil {
		t.Fatal(err)
	}
}
