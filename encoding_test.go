package utf8reader

import (
	"testing"
)

func TestTrunc(t *testing.T) {
	data := []struct {
		in  []byte
		out int
	}{
		// fully encoded "mouillé"
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xc3, 0xa9}, 8},
		// truncated "mouillé"
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xc3}, 6},
		// non truncated "=🏖️"
		{[]byte{0x3d, 0xf0, 0x9f, 0x8f, 0x96, 0xef, 0xb8, 0x8f}, 8},
		// truncated "=🏖️"
		{[]byte{0x3d, 0xf0, 0x9f, 0x8f, 0x96, 0xef, 0xb8}, 5},
		{[]byte{0x3d, 0xf0, 0x9f, 0x8f, 0x96, 0xef}, 5},
		{[]byte{0x3d, 0xf0, 0x9f, 0x8f, 0x96}, 5},
		{[]byte{0x3d, 0xf0, 0x9f, 0x8f}, 1},
		{[]byte{0x3d, 0xf0, 0x9f}, 1},
		{[]byte{0x3d, 0xf0}, 1},
		{[]byte{0x3d}, 1},
	}
	for _, d := range data {
		if got := trunc(d.in); got != d.out {
			t.Errorf("trunc(%v) = %v, want %v", d.in, got, d.out)
		}
	}
}

func TestIsUTF8(t *testing.T) {
	data := []struct {
		in  []byte
		out bool
	}{
		// fully encoded "mouillé"
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xc3, 0xa9}, true},
		// truncated "mouillé"
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xc3}, true},
		// iso8859-1 "mouillé" is a false positive because it is truncated "mouill�"
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xe9}, true},
		// iso8859-1 "mouillé "
		{[]byte{0x6d, 0x6f, 0x75, 0x69, 0x6c, 0x6c, 0xe9, 0x20}, false},
		// kio8-r "тест"
		{[]byte{0xf4, 0xe5, 0xf1, 0xf2}, false},
	}
	for _, d := range data {
		if got := isUTF8(d.in); got != d.out {
			t.Errorf("isUTF8(%v) = %v, want %v", d.in, got, d.out)
		}
	}
}

func TestGuessUTF16(t *testing.T) {
	data := []struct {
		in  []byte
		out string
	}{
		// UTF-16BE with BOM
		{[]byte{0xfe, 0xff, 0x00, 0x61, 0x00, 0x62}, "UTF-16BE"},
		// UTF-16LE with BOM
		{[]byte{0xff, 0xfe, 0x61, 0x00, 0x62, 0x00}, "UTF-16LE"},
		// UTF-16BE without BOM
		{[]byte{0x00, 0x61, 0x00, 0x62}, "UTF-16BE"},
		// UTF-16LE without BOM
		{[]byte{0x61, 0x00, 0x62, 0x00}, "UTF-16LE"},
		// UTF-16BE without BOM truncated
		{[]byte{0x00, 0x61, 0x00}, "UTF-16BE"},
		// UTF-16LE without BOM truncated
		{[]byte{0x61, 0x00, 0x62}, "UTF-16LE"},
	}
	for _, d := range data {
		if got := guessUTF16(d.in); got != d.out {
			t.Errorf("guessUTF16(%v) = %v, want %v", d.in, got, d.out)
		}
	}
}
