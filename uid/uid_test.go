package uid

import (
	"fmt"
	"testing"
)

func Test_SnowflakeID(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(SnowflakeID())
	}
}

func Test_SnowflakeIDString(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(SnowflakeIDString())
	}
}

func Test_SnowflakeIDFrom(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(SnowflakeIDFrom(_uuid.rand.Uint64()))
	}
}

func Test_UUID(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(UUID(true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID(true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID(false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID(false, false))
	}
}

func Test_UUID1(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(UUID1(true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID1(true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID1(false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID1(false, false))
	}
}

func Test_UUID2(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(UUID2(true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID2(true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID2(false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID2(false, false))
	}
}

func Test_UUID3(t *testing.T) {
	namespace, name := []byte("uuid"), []byte("v3")
	for i := 0; i < 5; i++ {
		fmt.Println(UUID3(namespace, name, true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID3(namespace, name, true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID3(namespace, name, false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID3(namespace, name, false, false))
	}
}

func Test_UUID4(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(UUID4(true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID4(true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID4(false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID4(false, false))
	}
}

func Test_UUID5(t *testing.T) {
	namespace, name := []byte("uuid"), []byte("v5")
	for i := 0; i < 5; i++ {
		fmt.Println(UUID5(namespace, name, true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID5(namespace, name, true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID5(namespace, name, false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUID5(namespace, name, false, false))
	}
}

func Test_UUIDFrom(t *testing.T) {
	n1, n2 := _uuid.rand.Uint64(), _uuid.rand.Uint64()
	for i := 0; i < 5; i++ {
		fmt.Println(UUIDFrom(n1, n2, true, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUIDFrom(n1, n2, true, false))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUIDFrom(n1, n2, false, true))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(UUIDFrom(n1, n2, false, false))
	}
}
