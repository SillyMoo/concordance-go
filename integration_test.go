package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestTheLot(t *testing.T) {
	fTestText, err := os.Open("./test.text")
	if err != nil {
		t.Errorf("Could not open test text file, %v", err)
	}

	var fTestExpect *os.File
	fTestExpect, err = os.Open("./test.expect")
	rTestExpect := bufio.NewReader(fTestExpect)

	var out = new(bytes.Buffer)
	produceConcordance(fTestText, out, standardConcordanceLineOutput)

	var err2 error
	var str, str2 string

	for str, err = out.ReadString('\n'); err == nil; str, err = out.ReadString('\n') {
		str2, err2 = rTestExpect.ReadString('\n')
		fmt.Printf("%s %s\n", str, str2)
		if err2 == io.EOF {
			t.Errorf("Too many lines returned")
			return
		}
		if err2 != nil {
			t.Errorf("Unexpected error, %v", err2)
			return
		}

		if str != str2 {
			t.Errorf("Expected %s to equal %s", str, str2)
			return
		}
	}

	if err != io.EOF {
		t.Errorf("Unexpected error, %v", err)
	}

	_, err2 = rTestExpect.ReadString('\n')
	if err2 != io.EOF {
		t.Errorf("Not enough lines returned")
	}
}
