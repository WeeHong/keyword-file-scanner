package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	ignored := make(map[string]struct{})
	ignored[".git"] = struct{}{}
	p, _ := scanDirectory(ignored)
	check := false
	for _, path := range p {
		if strings.Contains(path, ".git") {
			check = true
		}
	}

	if check {
		t.Logf(".git path has been exlcuded from the search")
	} else {
		t.Errorf(".git path should not be found")
	}
}

func TestIsDetectKeywordWithAbsolutePath(t *testing.T) {
	isFound := false
	isDetectKeyword(true, "./main.go", "TODO", &isFound)

	if !isFound {
		t.Logf("Keyword \"TODO\" should not be found")
	} else {
		t.Errorf("Keyword \"TODO\" detected in the file")
	}
}

func TestIsDetectKeywordWithoutAbsolutePath(t *testing.T) {
	isFound := false
	isDetectKeyword(false, "./main.go", "TODO", &isFound)

	if !isFound {
		t.Logf("Keyword \"TODO\" should not be found")
	} else {
		t.Errorf("Keyword \"TODO\" detected in the file")
	}
}

func TestShowLineContainsText(t *testing.T) {
	buffer := bytes.Buffer{}
	ShowLineContainsText(&buffer, "./main.go", "TODO")

	got := buffer.String()
	want := "Line 163"

	if got != want {
		t.Logf("Keyword \"TODO\" should not be found")
	} else {
		t.Errorf("Keyword \"TODO\" found at: %s", want)
	}
}
