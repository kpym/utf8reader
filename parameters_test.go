package utf8reader

import (
	"golang.org/x/text/unicode/norm"
	"testing"
)

func TestNewParams(t *testing.T) {
	p := newParams()
	if p.peekSize != 4096 {
		t.Errorf("newParams().peakSize = %v, want 4096", p.peekSize)
	}
	if p.transformers != nil {
		t.Errorf("newParams().transformers = %v, want nil", p.transformers)
	}
	p = newParams(WithPeekSize(8192))
	if p.peekSize != 8192 {
		t.Errorf("newParams(WithPeakSize(8192)).peakSize = %v, want 8192", p.peekSize)
	}
	if p.transformers != nil {
		t.Errorf("newParams(WithPeakSize(8192)).transformers = %v, want nil", p.transformers)
	}
	p = newParams(WithNormalization("NFC"))
	if p.peekSize != 4096 {
		t.Errorf("newParams(WithNormalizationForm(\"NFC\")).peakSize = %v, want 4096", p.peekSize)
	}
	if len(p.transformers) != 1 || p.transformers[0] != norm.NFC {
		t.Errorf("newParams(WithNormalizationForm(\"NFC\")).transformers = %v, want [NFC]", p.transformers)
	}
	p = newParams(WithNormalization("NFD"))
	if p.peekSize != 4096 {
		t.Errorf("newParams(WithNormalizationForm(\"NFD\")).peakSize = %v, want 4096", p.peekSize)
	}
	if len(p.transformers) != 1 || p.transformers[0] != norm.NFD {
		t.Errorf("newParams(WithNormalizationForm(\"NFD\")).transformers = %v, want [NFD]", p.transformers)
	}
	p = newParams(WithPeekSize(8192), WithNormalization("NFC"))
	if p.peekSize != 8192 {
		t.Errorf("newParams(WithPeakSize(8192), WithNormalizationForm(\"NFC\")).peakSize = %v, want 8192", p.peekSize)
	}
	if len(p.transformers) != 1 || p.transformers[0] != norm.NFC {
		t.Errorf("newParams(WithPeakSize(8192), WithNormalizationForm(\"NFC\")).transformers = %v, want [NFC]", p.transformers)
	}
}
