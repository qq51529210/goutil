package discovery

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	ms, err := DiscoverMulticast(context.Background(), "en0", MulticastAddr, time.Second*5)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range ms {
		fmt.Println(m)
	}
	fmt.Println("1")
	ms, err = DiscoverAddress(context.Background(), "192.168.31.3:3702", time.Second*5)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range ms {
		fmt.Println(m)
	}
	fmt.Println("2")
	ms, err = DiscoverAddress(context.Background(), "192.168.31.255:3702", time.Second*5)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range ms {
		fmt.Println(m)
	}
	fmt.Println("3")
}
