package main

import (
	"fmt"
	"github.com/SillyMoo/concordance-go/tokeniser"
	"io"
	"sort"
	"strconv"
	"strings"
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
func generateWordPositions(chIn chan string, chOut chan wordPosition) {
	sentenceIdx := 0
	for str := range chIn {
		wordChan := make(chan string, 5)
		go func() {
			tokeniser.TokenFactory(tokeniser.ScanWords)(strings.NewReader(str), wordChan)
			close(wordChan)
		}()
		wordIdx := 0
		for str := range wordChan {
			chOut <- wordPosition{position{sentenceIdx, wordIdx}, strings.ToLower(str)}
			wordIdx++
		}
		sentenceIdx++
	}
	close(chOut)
}

//Given a channel of word positions, will produce a map of word to position array (i.e the concordance)
func generateConcordance(chIn chan wordPosition) (res map[string][]position) {
	res = make(map[string][]position)
	for pos := range chIn {
		res[pos.word] = append(res[pos.word], position{pos.sentenceIdx, pos.wordIdx})
	}
	return
}

//Used for outputing lines of the concordance file
type concordanceLineOutput func(word string, positions []position, r io.Writer)

//Standard format concordance line outputter
func standardConcordanceLineOutput(word string, positions []position, r io.Writer) {
	var strArr []string
	for _, pos := range positions {
		strArr = append(strArr, strconv.Itoa(pos.sentenceIdx))
	}
	fmt.Fprintf(r, "%s {%d:%v}\n", word, len(positions), strings.Join(strArr, ","))
}

//Concordance line outputter that will include word position as well as the usual frequencey and sentence position
func wordPositionConcordanceLineOutput(word string, positions []position, r io.Writer) {
	var strArr []string
	for _, pos := range positions {
		strArr = append(strArr, strconv.Itoa(pos.sentenceIdx)+"."+strconv.Itoa(pos.wordIdx))
	}
	fmt.Fprintf(r, "%s {%d:%v}\n", word, len(positions), strings.Join(strArr, ","))
}

//Outputs a concordance to a reader. Takes a function that will provide the desired output for a single concordance line.
//Takes care of sorting the output by alphanumeric order.
func outputConcordance(r io.Writer, res map[string][]position, c concordanceLineOutput) {
	var keys = make([]string, len(res))
	i := 0
	for k, _ := range res {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		c(k, res[k], r)
	}
}
