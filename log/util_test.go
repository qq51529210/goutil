package log

import "testing"

func Test_ParseSize(t *testing.T) {
	for _, s := range []struct {
		s string
		n float64
	}{
		{"13.5TB", 13.5 * tb},
		{"1.5T", 1.5 * tb},
		{"0.3GB", 0.3 * gb},
		{"0.31G", 0.31 * gb},
		{"12MB", 12 * mb},
		{"2.132M", 2.132 * mb},
		{"30KB", 30 * kb},
		{"0.34K", 0.34 * kb},
		{"123", 123},
	} {
		n, err := ParseSize(s.s)
		if err != nil {
			t.Fatal(err)
		}
		if n != int64(s.n) {
			t.FailNow()
		}
	}
}
