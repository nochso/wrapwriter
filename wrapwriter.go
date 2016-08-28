package wrapwriter

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

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

func nextEOW(data []byte, offset int) (end int) {
	end = bytes.IndexByte(data[offset:], ' ')
	if end == -1 {
		end = len(data)
	} else {
		end += offset
	}
	return
}

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
				out.WriteByte(' ')
				remaining--
			}
			// fill rest of line
			out.WriteString(word[0:remaining])
			out.WriteByte('\n')
			wStart := remaining
			remaining = width
			for wStart < wordLen {
				if remaining <= 0 {
					out.WriteByte('\n')
					remaining = width
				}
				wEnd := wStart + width
				if wEnd > wordLen {
					wEnd = wordLen
				}
				out.WriteString(word[wStart:wEnd])
				remaining = width - (wEnd - wStart)
				wStart = wEnd
			}
		}
		firstWord = false
		offset = eow + 1
	}
	if hasEOL {
		out.WriteByte('\n')
	}
}
