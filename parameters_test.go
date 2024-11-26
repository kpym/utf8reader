package utf8reader

import (
	"golang.org/x/text/unicode/norm"
	"testing"
)

func TestNewParams(t *testing.T) {
	p := newParams()
	if p.peakSize != 4096 {
		t.Errorf("newParams().peakSize = %v, want 4096", p.peakSize)
	}
	if p.norm != nil {
		t.Errorf("newParams().norm = %v, want nil", p.norm)
	}
	p = newParams(WithPeakSize(8192))
	if p.peakSize != 8192 {
		t.Errorf("newParams(WithPeakSize(8192)).peakSize = %v, want 8192", p.peakSize)
	}
	if p.norm != nil {
		t.Errorf("newParams(WithPeakSize(8192)).norm = %v, want nil", p.norm)
	}
	p = newParams(WithNormalizationForm("NFC"))
	if p.peakSize != 4096 {
		t.Errorf("newParams(WithNormalizationForm(\"NFC\")).peakSize = %v, want 4096", p.peakSize)
	}
	if p.norm != norm.NFC {
		t.Errorf("newParams(WithNormalizationForm(\"NFC\")).norm = %v, want NFC", p.norm)
	}
	p = newParams(WithNormalizationForm("NFD"))
	if p.peakSize != 4096 {
		t.Errorf("newParams(WithNormalizationForm(\"NFD\")).peakSize = %v, want 4096", p.peakSize)
	}
	if p.norm != norm.NFD {
		t.Errorf("newParams(WithNormalizationForm(\"NFD\")).norm = %v, want NFD", p.norm)
	}
	p = newParams(WithPeakSize(8192), WithNormalizationForm("NFC"))
	if p.peakSize != 8192 {
		t.Errorf("newParams(WithPeakSize(8192), WithNormalizationForm(\"NFC\")).peakSize = %v, want 8192", p.peakSize)
	}
	if p.norm != norm.NFC {
		t.Errorf("newParams(WithPeakSize(8192), WithNormalizationForm(\"NFC\")).norm = %v, want NFC", p.norm)
	}
}
