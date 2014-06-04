package main

import (
	"bufio"
	"flag"
	"github.com/SillyMoo/concordance-go/tokeniser"
	"io"
	"os"
)

//The main functionality for producing a concordance, pulled out in to a separate function
//to facilitate an 'integration test'
func produceConcordance(in io.Reader, out io.Writer, c concordanceLineOutput) {
	var ch1 = make(chan string)
	var ch2 = make(chan wordPosition)

	//Split incoming text in to sentences
	go func() {
		tokeniser.TokenFactory(bufio.ScanLines).Compose(
			tokeniser.TokenFactory(tokeniser.ScanSentences))(
			in, ch1)
		close(ch1)
	}()

	go generateWordPositions(ch1, ch2)

	outputConcordance(out, generateConcordance(ch2), c)
}

func main() {
	var format string
	flag.StringVar(&format, "outputFormat", "standard", "output format: [standard|wordIdx]")
	flag.Parse()
	outputFunction := standardConcordanceLineOutput
	if format == "wordIdx" {
		outputFunction = wordPositionConcordanceLineOutput
	}
	produceConcordance(os.Stdin, os.Stdout, outputFunction)
}
