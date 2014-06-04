package tokeniser

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

type Tokeniser func(r io.Reader, ch chan string)

//Simple factory function that produces a tokeniser given a splitting function
func TokenFactory(sf bufio.SplitFunc) Tokeniser {
	return func(r io.Reader, ch chan string) {
		var scanner = bufio.NewScanner(r)
		scanner.Split(sf)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}
}

//Splitting function that breaks a paragraph in to a set of sentences. Will generally split on full stops or exclamation
//marks, however it is aware of opening and closing brackets, quotes, etc.
//If a punctuation appears inside open/closing brackets it will not be counted as a sentence break. This is not too smart though
//and won't cope with nested brackets, etc.
func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//remove preceeding whitespace
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}
	enclosed := false

	//look for a period followed by whitespace
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !enclosed && unicode.In(r, unicode.Ps) {
			//an open bracket/quote, etc was found. Any punctuation found until the matching close will be ignored.
			enclosed = true
		} else if enclosed && unicode.In(r, unicode.Pe) {
			enclosed = false
		}
		if !enclosed && (r == '.' || r == '!') {
			if i+width < len(data) {
				r2, width2 := utf8.DecodeRune(data[i+width:])
				if unicode.IsSpace(r2) {
					return i + width + width2, data[start:i], nil
				}
			}
		}
	}
	//if hit end of file this must be a sentence (user just didn't add a full stop)
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return 0, nil, nil
}

//Splitting function that separates a sentence in to a set of words. Will remove whitespace, and will also remove punctuation at
//the end of a word (but will consider punctuation within a word as part of the word, such as 'e.g').
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	//remove preceeding whitespace and punctuation
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			break
		}
	}

	firstPunct := -1

	//look for next space character
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			//Found punctuation at end of a word, strip it
			if firstPunct != -1 {
				return i, data[start:firstPunct], nil
			} else {
				return i, data[start:i], nil
			}
		} else if unicode.IsPunct(r) && firstPunct == -1 {
			//Mark first punctuation found, this is used to strip punctuation within at the end of a word
			firstPunct = i
		} else if firstPunct != -1 && !unicode.IsPunct(r) {
			//Punctuation followed by non-punctuation is punctuation inside a word, we don't wish to strip that
			firstPunct = -1
		}
	}

	if atEOF && len(data) > start {
		if firstPunct != -1 {
			return len(data), data[start:firstPunct], nil
		} else {
			return len(data), data[start:len(data)], nil
		}
	}
	return 0, nil, nil
}
