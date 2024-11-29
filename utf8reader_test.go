package utf8reader

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"golang.org/x/text/unicode/norm"
)

func TestNew_nil(t *testing.T) {
	r := New(nil)
	if r != nil {
		t.Errorf("New(nil) = %v, want nil", r)
	}
	if r.Encoding() != "" {
		t.Errorf("New(nil).Encoding() = %s, want \"\"", r.Encoding())
	}
	if b, e := r.Peek(); b != nil || e != io.EOF {
		t.Errorf("New(nil).Peek() = %v, %v, want nil, io.EOF", b, e)
	}
	b := make([]byte, 1024)
	if n, e := r.Read(b); n != 0 || e != io.EOF {
		t.Errorf("New(nil).Read(nil) = %d, %v, want 0, io.EOF", n, e)
	}
}

func TestNew(t *testing.T) {
	sr := strings.NewReader("bête")
	r := New(sr)
	if r == nil {
		t.Errorf("New(strings.NewReader(\"bête\")) = nil, want *Reader")
	} else {
		if r.tr == nil {
			t.Errorf("New(strings.NewReader(\"bête\")).r = nil, non nil expected")
		}
		if r.t != nil {
			t.Errorf("New(strings.NewReader(\"bête\")).t = %v, want nil", r.t)
		}
		if string(r.buf) != "bête" {
			t.Errorf("New(strings.NewReader(\"bête\")).buf = %s, want \"bête\"", r.buf)
		}
	}

	sr.Seek(0, 0)
	r = New(sr, WithPeekSize(8192))
	if r == nil {
		t.Errorf("New(strings.NewReader(\"bête\"), WithPeekSize(8192)) = nil, want *Reader")
	} else {
		if r.tr == nil {
			t.Errorf("New(strings.NewReader(\"bête\"), WithPeekSize(8192)).r = nil, non nil expected")
		}
		if r.t != nil {
			t.Errorf("New(strings.NewReader(\"bête\"), WithPeekSize(8192)).t = %v, want nil", r.t)
		}
		if string(r.buf) != "bête" {
			t.Errorf("New(strings.NewReader(\"bête\"), WithPeekSize(8192)).buf = %s, want \"bête\"", r.buf)
		}
	}

	sr.Seek(0, 0)
	r = New(sr, WithNormalizationForm("NFD"))
	if r == nil {
		t.Errorf("New(strings.NewReader(\"bête\"), WithNormalizationForm(\"NFD\")) = nil, want *Reader")
	} else {
		if r.tr == nil {
			t.Errorf("New(strings.NewReader(\"bête\"), WithNormalizationForm(\"NFD\")).r = nil, non nil expected")
		}
		if r.t != norm.NFD {
			t.Errorf("New(strings.NewReader(\"bête\"), WithNormalizationForm(\"NFD\")).t = %v, want %v", r.t, norm.NFD)
		}
		// check if the buffer contains NFD representation of "bête"
		if string(r.buf) != "bête" {
			t.Errorf("New(strings.NewReader(\"bête\"), WithPeekSize(8192)).buf = %s, want \"bête\"", r.buf)
		}
	}
}

func TestPeek(t *testing.T) {
	r := New(strings.NewReader("test"))
	if r == nil {
		t.Errorf("New(strings.NewReader(\"test\")) = nil, want *Reader")
	}
	s, err := r.Peek()
	if err != nil {
		t.Errorf("r.Peek() = %v, want nil", err)
	}
	if string(s) != "test" {
		t.Errorf("r.Peek() = %s, want \"test\"", s)
	}

	// UTF-16LE without BOM content : "bétà"
	r = New(bytes.NewReader([]byte{0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0, 0x00}))
	if r == nil {
		t.Errorf("New(bytes.NewReader([]byte{0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0, 0x00})) = nil, want *Reader")
	}
	s, err = r.Peek()
	if err != nil {
		t.Errorf("r.Peek() = %v, want nil", err)
	}
	if string(s) != "bétà" {
		t.Errorf("r.Peek() = %s, want \"bétà\"", s)
	}

	// UTF-16LE without BOM content, truncated from "bétà"
	r = New(bytes.NewReader([]byte{0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0}))
	if r == nil {
		t.Errorf("New(bytes.NewReader([]byte{0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0})) = nil, want *Reader")
	}
	s, err = r.Peek()
	if err != nil {
		t.Errorf("r.Peek() = %v, want nil", err)
	}
	if string(s) != "bét�" {
		t.Errorf("r.Peek() = %s, want \"bét�\"", s)
	}
}

func TestRead(t *testing.T) {
	data := []struct {
		name     string
		encoding string
		in       []byte
		out      []byte
	}{
		{
			name:     "bétà : UTF-8",
			encoding: "UTF-8",
			in:       []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
			out:      []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
		},
		{
			name:     "bétà : UTF-16LE",
			encoding: "UTF-16LE",
			in:       []byte{0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0, 0x00},
			out:      []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
		},
		{
			name:     "bétà : UTF-16LE with BOM",
			encoding: "UTF-16LE",
			in:       []byte{0xff, 0xfe, 0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0, 0x00},
			out:      []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
		},
		{
			name:     "bétà : UTF-16BE",
			encoding: "UTF-16BE",
			in:       []byte{0x00, 0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0},
			out:      []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
		},
		{
			name:     "bétà : UTF-16BE with BOM",
			encoding: "UTF-16BE",
			in:       []byte{0xfe, 0xff, 0x00, 0x62, 0x00, 0xe9, 0x00, 0x74, 0x00, 0xe0},
			out:      []byte{0x62, 0xc3, 0xa9, 0x74, 0xc3, 0xa0},
		},
		{
			name:     "C'est bête en français : iso-8859-1",
			encoding: "ISO-8859-1",
			in:       []byte{0x43, 0x27, 0x65, 0x73, 0x74, 0x20, 0x62, 0xEA, 0x74, 0x65, 0x20, 0x65, 0x6E, 0x20, 0x66, 0x72, 0x61, 0x6E, 0xE7, 0x61, 0x69, 0x73},
			out:      []byte{0x43, 0x27, 0x65, 0x73, 0x74, 0x20, 0x62, 0xC3, 0xAA, 0x74, 0x65, 0x20, 0x65, 0x6E, 0x20, 0x66, 0x72, 0x61, 0x6E, 0xC3, 0xA7, 0x61, 0x69, 0x73},
		},
		{
			name:     "Глупаво е на български : windows-1251",
			encoding: "WINDOWS-1251",
			in:       []byte{0xC3, 0xEB, 0xF3, 0xEF, 0xE0, 0xE2, 0xEE, 0x20, 0xE5, 0x20, 0xED, 0xE0, 0x20, 0xE1, 0xFA, 0xEB, 0xE3, 0xE0, 0xF0, 0xF1, 0xEA, 0xE8},
			out:      []byte{0xD0, 0x93, 0xD0, 0xBB, 0xD1, 0x83, 0xD0, 0xBF, 0xD0, 0xB0, 0xD0, 0xB2, 0xD0, 0xBE, 0x20, 0xD0, 0xB5, 0x20, 0xD0, 0xBD, 0xD0, 0xB0, 0x20, 0xD0, 0xB1, 0xD1, 0x8A, 0xD0, 0xBB, 0xD0, 0xB3, 0xD0, 0xB0, 0xD1, 0x80, 0xD1, 0x81, 0xD0, 0xBA, 0xD0, 0xB8},
		},
		{
			name:     "Глупаво е на български : koi8-r",
			encoding: "KOI8-R",
			in:       []byte{0xE7, 0xCC, 0xD5, 0xD0, 0xC1, 0xD7, 0xCF, 0x20, 0xC5, 0x20, 0xCE, 0xC1, 0x20, 0xC2, 0xDF, 0xCC, 0xC7, 0xC1, 0xD2, 0xD3, 0xCB, 0xC9},
			out:      []byte{0xD0, 0x93, 0xD0, 0xBB, 0xD1, 0x83, 0xD0, 0xBF, 0xD0, 0xB0, 0xD0, 0xB2, 0xD0, 0xBE, 0x20, 0xD0, 0xB5, 0x20, 0xD0, 0xBD, 0xD0, 0xB0, 0x20, 0xD0, 0xB1, 0xD1, 0x8A, 0xD0, 0xBB, 0xD0, 0xB3, 0xD0, 0xB0, 0xD1, 0x80, 0xD1, 0x81, 0xD0, 0xBA, 0xD0, 0xB8},
		},
	}

	buf := make([]byte, 1024)
	for _, d := range data {
		r := New(bytes.NewReader(d.in))
		if r == nil {
			t.Errorf("%s → New(bytes.NewReader(% X)) = nil, want *Reader", d.name, d.in)
		}
		if strings.ToUpper(r.Encoding()) != d.encoding {
			t.Errorf("%s → r.Encoding() = %s, want %s", d.name, r.Encoding(), d.encoding)
		}
		b := buf
		n, err := r.Read(b)
		if err != nil {
			t.Errorf("%s → r.Read(), err = %v, want nil", d.name, err)
		}
		b = b[:n]
		if !bytes.Equal(b, d.out) {
			t.Errorf("%s → r.Read(), b = % X, want % X", d.name, b, d.out)
		}
	}
}
