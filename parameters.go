package utf8reader

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// readerParams contains the parameters for the reader.
type readerParams struct {
	peakSize int                   // The number of bytes to peak
	norm     transform.Transformer // The normalization form NFC or NFD
}

// option is a functional option for the reader.
type option func(*readerParams)

// WithPeakSize sets the number of bytes to peak.
// By default it peaks 4096 bytes.
// The peaked bytes are used to detect the encoding.
func WithPeakSize(size int) option {
	return func(p *readerParams) {
		p.peakSize = size
	}
}

// WithNormalizationForm sets the normalization form.
// The normalization form can be "NFC" or "NFD".
// By default no normalization is done.
func WithNormalizationForm(nor string) option {
	return func(p *readerParams) {
		switch nor {
		case "NFC":
			p.norm = norm.NFC
		case "NFD":
			p.norm = norm.NFD
		default:
			panic("only NFC and NFD are supported")
		}
	}
}

// newParams returns a new readerParams with the options set.
func newParams(options ...option) *readerParams {
	p := &readerParams{
		peakSize: 4096,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}
