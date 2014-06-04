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
	ScannerTester(t, ScanWords, "a", []string{"a"})
	ScannerTester(t, ScanWords, "monkey magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanWords, "monkey, magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanWords, "monkey,, magic", []string{"monkey", "magic"})
	ScannerTester(t, ScanWords, "monkey, magic!", []string{"monkey", "magic"})
	ScannerTester(t, ScanWords, "monkey-magic!", []string{"monkey-magic"})
}

func TestDoubleByteWords(t *testing.T) {
	ScannerTester(t, ScanWords, "¢", []string{"¢"})
	ScannerTester(t, ScanWords, "a¢", []string{"a¢"})
	ScannerTester(t, ScanWords, "a¢, 1234", []string{"a¢", "1234"})
	ScannerTester(t, ScanWords, "a¢”", []string{"a¢"})
	ScannerTester(t, ScanWords, "a¢”, hello", []string{"a¢", "hello"})
	ScannerTester(t, ScanWords, "", []string{})
}

func TestSentence(t *testing.T) {
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle",
		[]string{"Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle"})
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle.",
		[]string{"Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle."})
	ScannerTester(t, ScanSentences, "Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle. Deep within this dank and uninviting place, lives Berk ",
		[]string{"Somewhere in the dark and nasty regions, where nobody goes, stands an ancient castle",
			"Deep within this dank and uninviting place, lives Berk "})
	ScannerTester(t, ScanSentences, "Don't you open that trapdoor, You're a fool if you dare.... Stay away from that trapdoor, 'Cos there's unicode down there ¢¢....",
		[]string{"Don't you open that trapdoor, You're a fool if you dare...",
			"Stay away from that trapdoor, 'Cos there's unicode down there ¢¢...."})
}

func TestSentenceWithBracketsAndQuotes(t *testing.T) {
	ScannerTester(t, ScanSentences, "Deep within this dank and uninviting place, lives Berk (Allo!), overworked servant of \"the thing upstairs\" (Berk! Feed Me!)",
		[]string{"Deep within this dank and uninviting place, lives Berk (Allo!), overworked servant of \"the thing upstairs\" (Berk! Feed Me!)"})
}

func TestTokenFactory(t *testing.T) {
	ch := make(chan string)
	go func() {
		TokenFactory(bufio.ScanRunes)(strings.NewReader("abcd"), ch)
		close(ch)
	}()
	var expect = []string{"a", "b", "c", "d"}
	i := 0
	for str := range ch {
		if str != expect[i] {
			t.Errorf("%s!=%s in %s", str, expect[i], "abcd")
		}
		i++
	}
	if i != len(expect) {
		t.Error("not enough tokens recieved from token factor")
	}
}

func TestTokeniserComposition(t *testing.T) {
	ch := make(chan string)
	go func() {
		TokenFactory(ScanWords).Compose(TokenFactory(bufio.ScanRunes))(strings.NewReader("ab cd"), ch)
		close(ch)
	}()
	var expect = []string{"a", "b", "c", "d"}
	i := 0
	for str := range ch {
		if str != expect[i] {
			t.Errorf("%s!=%s in %s", str, expect[i], "ab cd")
		}
		i++
	}
	if i != len(expect) {
		t.Error("not enough tokens recieved from token factory")
	}
}
