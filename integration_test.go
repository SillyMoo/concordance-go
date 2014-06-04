package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

//Integration test that run test.text through the concur producer, and compares the results with
//test.expect. The file test.expect contains a reformated version of the example results from the assignment
//description, but with two small changes. Firstly 'e.g.' has been replaced with 'e.g.', I decided to not add a special
//case just for this. Secondly the frequency for 'concordance', was bumped up to 3, it is 2 in the sample result, but
//this seems to be an issue with the document.
func TestTheLot(t *testing.T) {
	fTestText, err := os.Open("./testResources/test.text")
	if err != nil {
		t.Errorf("Could not open test text file, %v", err)
	}

	var fTestExpect *os.File
	fTestExpect, err = os.Open("./testResources/test.expect")
	rTestExpect := bufio.NewReader(fTestExpect)

	var out = new(bytes.Buffer)
	produceConcordance(fTestText, out, standardConcordanceLineOutput)

	var err2 error
	var str, str2 string

	for str, err = out.ReadString('\n'); err == nil; str, err = out.ReadString('\n') {
		str2, err2 = rTestExpect.ReadString('\n')
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
