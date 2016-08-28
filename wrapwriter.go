// Package wrapwriter wraps strings for console output.
package wrapwriter

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// Wrap text to width or fewer runes.
//
// Text with multiple lines is supported. Lines are assumed to use a single
// line feed "\n". Width must be positive.
func Wrap(text string, width int) (string, error) {
	if width <= 0 {
		return "", fmt.Errorf("expecting positive wrap width, got %d", width)
	}
	offset := 0
	out := bytes.NewBuffer(make([]byte, 0, int(float32(len(text))*1.05)))
	data := []byte(text)
	for offset < len(data) {
		// find and wrap a line
		eolPos, hasEOL := nextEOL(data, offset)
		wrapLine(out, data[offset:eolPos], hasEOL, width)
		offset = eolPos + 1
	}
	return out.String(), nil
}

// nextEOL looks for the end-of-line position. hasEOL may return false on the very last line.
func nextEOL(data []byte, offset int) (pos int, hasEOL bool) {
	pos = bytes.IndexByte(data[offset:], '\n')
	if pos == -1 {
		pos = len(data)
	} else {
		pos += offset
		hasEOL = true
	}
	return
}

// nextEOW looks for the end of the current word.
func nextEOW(data []byte, offset int) (end int) {
	end = bytes.IndexByte(data[offset:], ' ')
	if end == -1 {
		end = len(data)
	} else {
		end += offset
	}
	return
}

// wrapLine wraps a single line into buf.
func wrapLine(out *bytes.Buffer, line []byte, hasEOL bool, width int) {
	offset := 0
	firstWord := true
	remaining := width
	for offset < len(line) {
		if remaining <= 0 { // begin new line if previous line was full
			remaining = width
			out.WriteByte('\n')
			firstWord = true
		}
		// find end of word
		eow := nextEOW(line, offset)
		word := string(line[offset:eow])
		wordLen := utf8.RuneCountInString(word)
		if wordLen == 0 { // ignore leading spaces
			offset = eow + 1
			continue
		}
		sameLen := wordLen
		if !firstWord {
			sameLen++ // consider leading space
		}
		if sameLen <= remaining { // fits on remaining line
			if !firstWord {
				out.WriteByte(' ')
			}
			out.WriteString(word)
			remaining -= sameLen
		} else if wordLen <= width { // fits on its own line
			out.WriteByte('\n')
			out.WriteString(word)
			firstWord = false
			remaining = width - wordLen
		} else { // hard-wrap
			if !firstWord { // consider leading space
				if remaining > 1 {
					out.WriteByte(' ')
				}
				remaining--
			}
			for _, char := range word {
				if remaining <= 0 {
					out.WriteByte('\n')
					remaining = width
				}
				out.WriteRune(char)
				remaining--
			}
		}
		firstWord = false
		offset = eow + 1
	}
	if hasEOL {
		out.WriteByte('\n')
	}
}
