package utf8reader

import (
	"io"
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
