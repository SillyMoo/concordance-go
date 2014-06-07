package tokeniser

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
	//"fmt"
)

// A tokeniser is a function takes an incoming set of bytes, and produces a set of string tokens on a passed in channel
type Tokeniser func(r io.Reader, ch chan string)

func (t1 Tokeniser) Compose(t2 Tokeniser) Tokeniser {
	return func(r io.Reader, chOut chan string) {
		chMid := make(chan string)
		go func() {
			t1(r, chMid)
			close(chMid)
		}()
		for str := range chMid {
			t2(strings.NewReader(str), chOut)
		}
	}
}

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
	enclosed := false
	lastSentenceEnd := -1
	foundNonPunctuation := false
	firstPunct := -1
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) && (!unicode.IsPunct(r) || r=='.' || r=='!' || r=='\r' || r=='\n') {
			utf8.EncodeRune(data[start:], unicode.ToLower(r))
			break
		}
		if !enclosed && unicode.In(r, unicode.Ps) {
			//an open bracket/quote, etc was found. Any punctuation found until the matching close will be ignored.
			enclosed = true
		} else if enclosed && unicode.In(r, unicode.Pe) {
			enclosed = false
		}
	}

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
		/*if r=='\r' || r=='\n' {
		 	if foundNonPunctuation {
				end := i
				if lastSentenceEnd < firstPunct && lastSentenceEnd!=-1 {
				end = lastSentenceEnd
				} else if firstPunct != -1 {
				end = firstPunct
				}
				return i-1, data[start:end], nil
			} else {
				fmt.Println("b")
				return i, []byte{'.'}, nil
			}
		}*/
		if lastSentenceEnd == -1 && !enclosed && (r == '.' || r == '!') {
			lastSentenceEnd = i
			if firstPunct==-1 {
				firstPunct = i
			}
		} else if lastSentenceEnd == -1 {
			utf8.EncodeRune(data[i:], unicode.ToLower(r))
			foundNonPunctuation=true
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
		} else if lastSentenceEnd != -1 {
			if unicode.IsSpace(r) {
				if !foundNonPunctuation {
					return i, []byte{'.'}, nil
				}
				end := firstPunct
				if lastSentenceEnd < firstPunct {
					end = lastSentenceEnd
				}
				return i-1, data[start:end], nil
			} else if r != '.' && r != '!'{
				lastSentenceEnd=-1
			}
		}
	}
	//if hit end of file this must be a sentence (user just didn't add a full stop)
	if atEOF && len(data) > start {
		if lastSentenceEnd!=-1 {
			return len(data), data[start:lastSentenceEnd], nil
		} else if firstPunct!=-1 {
			return len(data), data[start:firstPunct], nil
		} else {
			return len(data), data[start:], nil
		}
	}
	return 0, nil, nil
}
