package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"strconv"
	"github.com/SillyMoo/concordance-go/tokeniser"
)

type wordPosition struct {
	sentenceIdx int
	wordIdx int
	word string
}

func main(){
	var ch1, ch2 = make(chan string), make(chan string)
	var ch3 = make(chan wordPosition)
	go func(){
		tokeniser.TokenFactory(bufio.ScanLines)(os.Stdin, ch1)
		close(ch1)
	}()
	go func(chIn, chOut chan string) {
		for str := range chIn {
			tokeniser.TokenFactory(tokeniser.ScanSentences)(strings.NewReader(str), chOut)
		}
		close(chOut)
	}(ch1, ch2)

	go func(chIn chan string, chOut chan wordPosition) {
		sentenceIdx :=0
		for str:= range chIn {
			var wordChan = make(chan string)
			go func(){
				tokeniser.TokenFactory(tokeniser.ScanWords)(strings.NewReader(str), wordChan)
				close(wordChan)
			}()
			wordIdx :=0
			for str := range wordChan {
				chOut <- wordPosition{sentenceIdx, wordIdx, strings.ToLower(str)}
				wordIdx++
			}
			sentenceIdx++
		}
		close(chOut)
	}(ch2, ch3)

	var res = make(map[string][]int)
	for pos := range ch3 {
		res[pos.word]= append(res[pos.word], pos.sentenceIdx)
	}
	for k,v := range res {
		var strArr []string
		for _, i := range v {
			strArr = append(strArr, strconv.Itoa(i))
		}
		fmt.Printf("%s {%d:%v}\n", k, len(v), strings.Join(strArr, ","))
	}
}