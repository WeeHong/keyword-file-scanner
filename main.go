package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	scanDirectory()
}

func scanDirectory() {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func filterDirectory(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Unable to read file from the path: %v", err)
	}
	if fi.Mode().IsRegular() {
		fmt.Println("file", path)
		f, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Unable to read file content: %v", err)
		}

		matched, err := regexp.Match(`Unable`, f)
		if err != nil {
			log.Fatalf("Unable to match the regexp pattern")
		}
		fmt.Println(matched)
	}
}

func showFilenameContainsText(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Unable to read file from the path: %v", err)
	}
	if fi.Mode().IsRegular() {
		fmt.Println("file", path)
		f, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Unable to read file content: %v", err)
		}

		matched, err := regexp.Match(`Unable`, f)
		if err != nil {
			log.Fatalf("Unable to match the regexp pattern")
		}
		fmt.Println(matched)
	}
}

func showLineContainsText(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open file")
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "Unable") {
			fmt.Println("Found at: ", line)
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
		fmt.Println(err)
	}
}
