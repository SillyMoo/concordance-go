package main

import (
	"bufio"
	"flag"
	"github.com/SillyMoo/concordance-go/tokeniser"
	"os"
)

func main() {
	var format string
	flag.StringVar(&format, "outputFormat", "standard", "output format: [standard|wordIdx]")
	flag.Parse()
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
	outputFunction := standardConcordanceLineOutput
	if format == "wordIdx" {
		outputFunction = wordPositionConcordanceLineOutput
	}
	outputConcordance(os.Stdout, generateConcordance(ch2), outputFunction)
}
