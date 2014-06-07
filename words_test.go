package main

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"strings"
)

func TestGenerateWordPositions(t *testing.T) {
	in := "He's the greatest. He's fantastic. Wherever there's danger he'll be there"
	chOut := make(chan wordPosition)

	go func() {
		generateWordPositions(strings.NewReader(in), chOut)
	}()
	numOfHes, foundBe := 0, false
	for pos := range chOut {
		switch pos.word {
		case "he's":
			if !(pos.sentenceIdx == 0 || pos.sentenceIdx == 1) {
				t.Error("Wrong sentence position for He's")
			}
			if pos.wordIdx != 0 {
				t.Error("Wrong word position for He's")
			}
			numOfHes++
		case "be":
			if pos.sentenceIdx != 2 {
				t.Error("Wrong sentence position for be")
			}
			if pos.wordIdx != 4 {
				t.Error("Wrong word position for be")
			}
			foundBe = true
		}
	}
	if numOfHes != 2 {
		t.Errorf("He's not found enough")
	}
	if !foundBe {
		t.Error("Didn't find he")
	}
}

func TestGenerateConcordance(t *testing.T) {
	in := []wordPosition{
		wordPosition{position{0, 0}, "the"},
		wordPosition{position{1, 0}, "the"},
		wordPosition{position{2, 0}, "the"},
		wordPosition{position{0, 1}, "pink"},
		wordPosition{position{0, 2}, "panther"}}

	ch := make(chan wordPosition)
	go func() {
		for _, wp := range in {
			ch <- wp
		}
		close(ch)
	}()
	res := generateConcordance(ch)
	if positions, ok := res["the"]; !ok {
		t.Error("'the' not found")
	} else {
		if len(positions) != 3 {
			t.Errorf("Expected there to be 3 occurences of 'the', but found %d", len(positions))
		}
	}

	if positions, ok := res["pink"]; !ok {
		t.Error("'pink' not found")
	} else {
		if len(positions) != 1 {
			t.Errorf("Expected there to be 1 occurences of 'pink', but found %d", len(positions))
		}
		if positions[0].sentenceIdx != 0 {
			t.Error("Expected pink sentenceIdx to be 0")
		}
		if positions[0].wordIdx != 1 {
			t.Error("Expected pink wordIdx to be 1")
		}
	}

	if positions, ok := res["panther"]; !ok {
		t.Error("'panther' not found")
	} else {
		if len(positions) != 1 {
			t.Errorf("Expected there to be 1 occurences of 'panther', but found %d", len(positions))
		}
		if positions[0].sentenceIdx != 0 {
			t.Error("Expected panther sentenceIdx to be 0")
		}
		if positions[0].wordIdx != 2 {
			t.Error("Expected panther wordIdx to be 2")
		}
	}
}

func TestOutputConcordance(t *testing.T) {
	in := []wordPosition{
		wordPosition{position{0, 0}, "the"},
		wordPosition{position{1, 0}, "the"},
		wordPosition{position{2, 0}, "the"},
		wordPosition{position{0, 1}, "pink"},
		wordPosition{position{0, 2}, "panther"}}

	expect := []string{
		"panther {1:0}",
		"pink {1:0}",
		"the {3:0,1,2}"}

	ch := make(chan wordPosition)
	go func() {
		for _, wp := range in {
			ch <- wp
		}
		close(ch)
	}()
	res := generateConcordance(ch)
	buffer := new(bytes.Buffer)
	inCh := bufio.NewWriter(buffer)
	outputConcordance(inCh, res, standardConcordanceLineOutput)

	for _, expected := range expect {
		line, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			} else {
				t.Error(err)
			}
		} else {
			if line != expected {
				t.Errorf("Expected: %s, got: %s", expected, line)
			}
		}
	}
}
