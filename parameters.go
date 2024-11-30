package utf8reader

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// readerParams contains the parameters for the reader.
type readerParams struct {
	peekSize     int                     // The number of bytes to peak
	transformers []transform.Transformer // The normalization form NFC or NFD
}

// option is a functional option for the reader.
type option func(*readerParams)

// WithPeekSize sets the number of bytes to peak.
// By default it peaks 4096 bytes.
// The peaked bytes are used to detect the encoding.
func WithPeekSize(size int) option {
	return func(p *readerParams) {
		p.peekSize = size
	}
}

// WithNormalization sets the normalization form.
// The normalization form can be "NFC" or "NFD".
// By default no normalization is done.
// WithNormalization("NFC") is equivalent to WithTransformers(norm.NFC).
// WithNormalization("NFD") is equivalent to WithTransformers(norm.NFD).
func WithNormalization(nor string) option {
	return func(p *readerParams) {
		switch nor {
		case "NFC":
			p.transformers = append(p.transformers, norm.NFC)
		case "NFD":
			p.transformers = append(p.transformers, norm.NFD)
		default:
			panic("only NFC and NFD are supported")
		}
	}
}

// WithTransformers append a (set of) transformer(s).
func WithTransform(transformers ...transform.Transformer) option {
	return func(p *readerParams) {
		p.transformers = append(p.transformers, transformers...)
	}
}

// newParams returns a new readerParams with the options set.
func newParams(options ...option) *readerParams {
	p := &readerParams{
		peekSize: 4096,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}
