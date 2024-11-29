// utf8reader is a package that detects the encoding of a reader
// and provides a new reader that converts (if necessary) the input to UTF-8 (without BOM).
// The unicode normalization form can be set to NFC or NFD.
package utf8reader

import (
	"io"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// peekReader allows to peek the first bytes of a reader.
// buf contains the first bytes of the reader.
// buf is set to nil when the buffer is empty.
// r is the underlying reader.
type peekReader struct {
	buf []byte
	r   io.Reader
}

// newPeekReader returns a new peekReader that peeks the first n bytes of the reader
// and stores them in the buffer.
// If some error occurs while reading the first n bytes, a nil peekReader is returned.
func newPeekReader(r io.Reader, n int) (*peekReader, error) {
	// no small buffer is allowed
	if n < 1024 {
		n = 1024
	}
	// create the peekReader
	pr := &peekReader{
		buf: make([]byte, n),
		r:   r,
	}
	// read the first n bytes
	n, err := io.ReadFull(r, pr.buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	pr.buf = pr.buf[:n]
	return pr, nil
}

// Read reads first from the peek buffer and then from the underlying reader.
func (r *peekReader) Read(p []byte) (n int, err error) {
	if r == nil {
		return 0, io.EOF
	}
	if len(r.buf) > 0 {
		n = copy(p, r.buf)
		r.buf = r.buf[n:]
		if len(r.buf) == 0 {
			// we need no buffer anymore, so it can be garbage collected
			r.buf = nil
		}
		return n, nil
	}
	return r.r.Read(p)
}

// peek returns the peek buffer.
// This function should be called before any Read operation.
func (r *peekReader) peek() []byte {
	return r.buf
}

// skip skips at most n bytes from the buffer.
func (r *peekReader) skip(n int) {
	if n > len(r.buf) {
		n = len(r.buf)
	}
	r.buf = r.buf[n:]
}

// Reader is a reader that converts the input to UTF-8.
type Reader struct {
	enc string                // the detected encoding
	buf []byte                // the peek buffer used to detect the encoding
	t   transform.Transformer // the encoding transformer & possibly the normalization transformer
	tr  io.Reader             // the underlying reader
}

// Read reads the UTF-8 encoded bytes from the reader.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r == nil {
		return 0, io.EOF
	}
	r.buf = nil
	return r.tr.Read(p)
}

// Peek returns the first bytes of the reader transformed to UTF-8.
// This function should be called before any Read operation.
func (r *Reader) Peek() ([]byte, error) {
	if r.buf == nil {
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

// Encoding returns the detected encoding.
func (r *Reader) Encoding() string {
	return r.enc
}

// New returns a reader that converts the input to UTF-8
// if it is not already encoded in UTF-8.
// If the encoding cannot be detected it returns buffered version of the original reader.
func New(r io.Reader, options ...option) *Reader {
	if r == nil {
		return &Reader{}
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
		if encoding != "UTF-8" {
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
	if len(trs) > 1 {
		tr = transform.Chain(trs...)
	} else if len(trs) == 1 {
		tr = trs[0]
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
