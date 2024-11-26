package utf8reader

import (
	"unicode/utf8"

	"github.com/gogs/chardet"
)

// trunc returns the length without possibly the last truncated rune.
func trunc(data []byte) int {
	end := len(data)
	lim := end - utf8.UTFMax
	if lim < 0 {
		lim = 0
	}
	for start := end - 1; start >= lim; start-- {
		if utf8.RuneStart(data[start]) {
			// try to decode the rune
			r, size := utf8.DecodeRune(data[start:])
			if r == utf8.RuneError {
				return start
			}
			return start + size
		}
	}
	return end
}

// isUTF8 returns true if the data is valid UTF-8,
// with possibly the last rune truncated.
func isUTF8(data []byte) bool {
	return utf8.Valid(data[:trunc(data)])
}

// guessUTF16 returns the "UTF-16 LE", "UTF-16 BE" if it looks like a valid UTF-16.
// - first it checks if the data starts with a BOM
// - if no bom is found it counts the number of
//   - <null><ascii> pairs (for UTF-16 BE)
//   - <ascii><null> pairs (for UTF-16 LE)
//
// Normally the other encodings do not have such pairs.
// We need this heuristic because chardet does not always detect UTF-16 correctly.
// For example, if the text is an ascii encoded as UTF-16 it will detect it as ASCII.
func guessUTF16(data []byte) string {
	if len(data) >= 2 {
		switch {
		case data[0] == 0xFE && data[1] == 0xFF:
			return "UTF-16BE"
		case data[0] == 0xFF && data[1] == 0xFE:
			return "UTF-16LE"
		}
	}
	utf16be := 0
	for i := 0; i < len(data)-1; i += 2 {
		if data[i] == 0 && data[i+1] < 128 {
			utf16be++
		}
	}
	utf16le := 0
	for i := 0; i < len(data)-1; i += 2 {
		if data[i] < 128 && data[i+1] == 0 {
			utf16le++
		}
	}
	if utf16be > 0 || utf16le > 0 {
		if utf16be > utf16le {
			return "UTF-16BE"
		} else {
			return "UTF-16LE"
		}
	}
	return ""
}

// detectCharset returns the encoding of the data.
// if the data is Ascii it returns an empty string.
func detectCharset(data []byte) string {
	if isUTF8(data) {
		return "UTF-8"
	}
	if encoding := guessUTF16(data); encoding != "" {
		return encoding
	}
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
	if err != nil {
		return ""
	}
	return result.Charset
}
