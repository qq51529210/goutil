package discovery

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	ms, err := Discover(context.Background(), "en0", MulticastAddr, time.Second*5)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range ms {
		fmt.Println(m)
	}
}
