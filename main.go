package main

import (
	"flag"
	"io"
	"os"
)

//The main functionality for producing a concordance, pulled out in to a separate function
//to facilitate an 'integration test'
func produceConcordance(in io.Reader, out io.Writer, c concordanceLineOutput) {
	var ch2 = make(chan wordPosition, 5)

	go generateWordPositions(in, ch2)

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
