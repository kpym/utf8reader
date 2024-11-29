// Package utf8reader provides a utility to wrap an io.Reader that contains text
// in an arbitrary encoding and produce an io.Reader that outputs UTF-8 encoded text.
// The package automatically detects the original encoding and converts the input to UTF-8.
// Additionally, it can normalize the text to a specified Unicode normalization form (NFC or NFD).
package utf8reader

import (
	"io"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// Reader wraps an io.Reader to convert its input to UTF-8 encoding, if required.
type Reader struct {
	enc string                // the detected encoding
	buf []byte                // the peek buffer used to detect the encoding
	t   transform.Transformer // the encoding transformer & possibly the normalization transformer
	tr  io.Reader             // the underlying reader
}

// Read reads data from the underlying reader, ensuring it is UTF-8 encoded.
// It returns the number of bytes read into p and any error encountered.
// If the Reader is nil, it returns 0 and io.EOF.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r == nil {
		return 0, io.EOF
	}
	r.buf = nil
	return r.tr.Read(p)
}

// Peek returns a UTF-8 encoded snapshot of the first bytes of the reader,
// primarily for encoding detection. The size of the snapshot is at most
// the size of the peek buffer, set by the PeekSize option.
// This method must be called before any Read operations.
func (r *Reader) Peek() ([]byte, error) {
	if r == nil || r.buf == nil {
		return nil, io.EOF
	}
	// transform the buffer
	if r.t == nil {
		return r.buf, nil
	}
	// transform the buffer
	tbuf, _, err := transform.Bytes(r.t, r.buf)
	// ignore ErrShortSrc, we transform what we can
	if err == transform.ErrShortSrc {
		err = nil
	}
	return tbuf, err
}

// Encoding returns the encoding detected from the input, or an empty string
// if detection was unsuccessful, or an error occurred during the detection.
func (r *Reader) Encoding() string {
	if r == nil {
		return ""
	}
	return r.enc
}

// New creates a Reader that converts the input to UTF-8.
// If encoding detection fails the input stays unchanged,
// and Encoding() will return an empty string.
func New(r io.Reader, options ...option) *Reader {
	if r == nil {
		return nil
	}
	params := newParams(options...)

	// peek the first bytes to detect the encoding
	pr, err := newPeekReader(r, params.peekSize)
	if err != nil {
		return nil
	}
	var encoding string
	var trs []transform.Transformer
	if beginning := pr.peek(); len(beginning) > 0 {
		if bom, lb := detectBOM(beginning); bom != "" {
			encoding = bom
			pr.skip(lb)
		} else {
			encoding = detectCharset(beginning)
		}
		if encoding != "UTF-8" && encoding != "" {
			if e, _ := charset.Lookup(encoding); e != nil {
				trs = append(trs, e.NewDecoder())
			}
		}
	}
	if params.norm != nil {
		trs = append(trs, params.norm)
	}
	// set the buffer
	reader := &Reader{
		enc: encoding,
		buf: pr.peek(),
	}
	// chain the transformers
	var tr transform.Transformer
	if encoding != "" {
		if len(trs) > 1 {
			tr = transform.Chain(trs...)
		} else if len(trs) == 1 {
			tr = trs[0]
		}
	}
	// install the transformer
	if tr == nil {
		reader.tr = pr
		return reader
	}
	reader.t = tr
	reader.tr = transform.NewReader(pr, tr)
	// ready to read
	return reader
}
