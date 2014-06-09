package main

import (
	"fmt"
	"github.com/SillyMoo/concordance-go-fast/tokeniser"
	"io"
	"sort"
	"strconv"
	"strings"
	"bufio"
)

//Defines the position of a word within a document, where sentenceIdx is the index of the sentence in which
//the word lies, and wordIdx is the index of the word within that sentence
type position struct {
	sentenceIdx int
	wordIdx     int
}

//A combination of a position, and a word
type wordPosition struct {
	position
	word string
}

//Given a channel with incoming sentences, will produce a set of word positions on the out channel
func generateWordPositions(r io.Reader, chOut chan wordPosition) {
	sentenceIdx := 0
	wordIdx := 0
	var scanner = bufio.NewScanner(r)
	scanner.Split(tokeniser.ScanSentences)
	for scanner.Scan() {
	 	str := scanner.Text()
		if str == "." {
			sentenceIdx++
			wordIdx=0
		} else {
			chOut <- wordPosition{position{sentenceIdx, wordIdx}, str}
			wordIdx++
		}
	}
	close(chOut)
}

//Given a channel of word positions, will produce a map of word to position array (i.e the concordance)
func generateConcordance(chIn chan wordPosition) (res map[string][]position) {
	res = make(map[string][]position)
	for pos := range chIn {
		res[pos.word] = append(res[pos.word], pos.position)
	}
	return
}

//Used for outputing lines of the concordance file
type concordanceLineOutput func(word string, positions []position, r io.Writer)

//Standard format concordance line output, in format "word {frequency:s1, s2, ...}"
func standardConcordanceLineOutput(word string, positions []position, r io.Writer) {
	l :=len(positions)-1
	// 11 = 32bit string in base 10 + comma
	var bytes []byte = make([]byte, 0,(l+1)*11)
	for i, pos := range positions {
		bytes = strconv.AppendInt(bytes, int64(pos.sentenceIdx), 10)
		if i<l {
			bytes = append(bytes, byte(','))
		}
	}
	fmt.Fprintf(r, "%s {%d:%s}\n", word, len(positions), string(bytes))
}

//Concordance line output that will include word position as well as the usual frequency and sentence position
func wordPositionConcordanceLineOutput(word string, positions []position, r io.Writer) {
	l:=len(positions)-1
	// 22 = 2*32 bit string in base 10 + '.' + ','
	var bytes []byte = make([]byte, 0, (l+1)*22)
	for i, pos := range positions {
		bytes = strconv.AppendInt(bytes, int64(pos.sentenceIdx), 10)
		bytes = append(bytes, byte('.'))
		bytes = strconv.AppendInt(bytes, int64(pos.wordIdx), 10)
		if i<l {
			bytes = append(bytes, byte(','))
		}
	}
	fmt.Fprintf(r, "%s {%d:%s}\n", word, len(positions), string(bytes))
}

//Outputs a concordance to a reader. Takes a function that will provide the desired output for a single concordance line.
//Takes care of sorting the output by alphanumeric order.
func outputConcordance(r io.Writer, res map[string][]position, c concordanceLineOutput) {
	var keys = make([]string, len(res))
	i := 0
	for k := range res {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		c(k, res[k], r)
	}
}
