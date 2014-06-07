package tokeniser

import (
	"bufio"
	"strings"
	"testing"
)

func ScannerTester(t *testing.T, sf bufio.SplitFunc, in string, expect []string) {
	var scanner = bufio.NewScanner(strings.NewReader(in))
	scanner.Split(sf)
	i := 0
	for scanner.Scan() {
		if scanner.Text() != expect[i] {
			t.Errorf("%s!=%s in %s", scanner.Text(), expect[i], in)
			return
		}
		i++
	}
}

func TestAsciiWords(t *testing.T) {
	ScannerTester(t, ScanSentences, "a", []string{"a"})
	ScannerTester(t, ScanSentences, "monkey magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanSentences, "monkey, magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanSentences, "monkey,, magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanSentences, "monkey, magic!", []string{"monkey", "magic"})
	ScannerTester(t, ScanSentences, "monkey-magic!", []string{"monkey-magic"})
}

func TestDoubleByteWords(t *testing.T) {
	ScannerTester(t, ScanSentences, "¢", []string{"¢"})
	ScannerTester(t, ScanSentences, "a¢", []string{"a¢"})
	ScannerTester(t, ScanSentences, "a¢, 1234", []string{"a¢", "1234"})
	ScannerTester(t, ScanSentences, "a¢”", []string{"a¢"})
	ScannerTester(t, ScanSentences, "a¢”, hello", []string{"a¢", "hello"})
	ScannerTester(t, ScanSentences, "", []string{})
}

func TestSentence(t *testing.T) {
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle",
		[]string{"somewhere", "in", "the", "dark", "and", "nasty", "regions", "where", "nobody", "goes", "stands",
		 "an", "ancient", "castle"})
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle.",
		[]string{"somewhere", "in", "the", "dark", "and", "nasty", "regions", "where", "nobody", "goes", "stands",
		 "an", "ancient", "castle"})
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle. Deep within this dank and uninviting place, lives Berk ",
		[]string{"somewhere", "in", "the", "dark", "and", "nasty", "regions", "where", "nobody", "goes", "stands",
		 "an", "ancient", "castle", ".", "deep", "within", "this", "dank", "and", "uninviting", "place", "lives", "berk"})
	ScannerTester(t, ScanSentences, "Don't you open that trapdoor, You're a fool if you dare.... Stay away from that trapdoor, 'Cos there's unicode down there ¢¢....",
		[]string{"don't", "you", "open", "that", "trapdoor", "you're", "a", "fool", "if", "you", "dare", ".",
		 "stay", "away", "from", "that", "trapdoor", "cos", "there's", "unicode", "down", "there", "¢¢"})
}

func TestSentenceWithBracketsAndQuotes(t *testing.T) {
	ScannerTester(t, ScanSentences, "Deep within this dank and uninviting place, lives Berk (Allo!), overworked servant of \"the thing upstairs\" (Berk! Feed Me!)",
		[]string{"deep", "within", "this", "dank", "and", "uninviting", "place", "lives", "berk", "allo", "overworked",
		 "servant", "of", "the", "thing", "upstairs", "berk", "feed", "me"})
}
