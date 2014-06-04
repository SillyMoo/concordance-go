package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"strconv"
	"github.com/SillyMoo/concordance-go/tokeniser"
)

//Defines the position of a word within a document, where sentenceIdx is the index of the sentence in which
//the word lies, and wordIdx is the index of the word within that sentence
type position struct {
	sentenceIdx int
	wordIdx int
}

//A combination of a position, and a word
type wordPosition struct {
	position
	word string
}

//Given a channel with incoming sentences, will produce a set of word positions on the out channel
func generateWordPositions(chIn chan string, chOut chan wordPosition) {
	sentenceIdx :=0
	for str:= range chIn {
		var wordChan = make(chan string)
		go func(){
			tokeniser.TokenFactory(tokeniser.ScanWords)(strings.NewReader(str), wordChan)
			close(wordChan)
		}()
		wordIdx :=0
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
		res[pos.word]= append(res[pos.word], position{pos.sentenceIdx, pos.wordIdx})
	}
	return
}

func main(){
	var ch1 = make(chan string)
	var ch2 = make(chan wordPosition)

	//Split incoming text in to sentences
	go func() {
		tokeniser.TokenFactory(bufio.ScanLines).Compose(
			tokeniser.TokenFactory(tokeniser.ScanSentences))(
				os.Stdin, ch1)
		close(ch1)
	}()

	go generateWordPositions(ch1, ch2)

	var res = generateConcordance(ch2)

	for k,v := range res {
		var strArr []string
		for _, pos := range v {
			strArr = append(strArr, strconv.Itoa(pos.sentenceIdx))
		}
		fmt.Printf("%s {%d:%v}\n", k, len(v), strings.Join(strArr, ","))
	}
}