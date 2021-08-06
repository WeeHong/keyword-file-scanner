package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type folder struct {
	List map[string]bool
}

type myCloser interface {
	Close() error
}

// closeFile is a helper function which streamlines closing
// with error checking on different file types.
func closeFile(f myCloser) {
	err := f.Close()
	check(err)
}

// readAll is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func readAll(file *zip.File) []byte {
	fc, err := file.Open()
	check(err)
	defer closeFile(fc)

	content, err := ioutil.ReadAll(fc)
	check(err)

	return content
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	collection := make(map[string]bool)

	for scanner.Scan() {
		key := strings.Trim(scanner.Text(), " ")
		if _, ok := collection[key]; !ok {
			collection[key] = true
		}
	}

	ignoredFolder.List = collection

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var ignoredFolder folder

func init() {
	if _, err := os.Stat(".folderignore"); err == nil {
		file, err := os.Open(".folderignore")
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()
		readFile(file)
	}
}

func main() {
	// zipFile := "./compress/aws-sdk-go-main.zip"
	// zf, err := zip.OpenReader(zipFile)
	// check(err)
	// defer closeFile(zf)

	// for _, file := range zf.File {
	// 	fmt.Printf("=%s\n", file.Name)
	// 	// fmt.Printf("%s\n\n", readAll(file)) // file content
	// }
	fmt.Println(ignoredFolder)
	_ = scanDirectory()
	// fmt.Println(filePaths)

	// keywords := os.Args[1:]
	// for _, keyword := range keywords {
	// 	fmt.Println(keyword)
	// }
}

func scanDirectory() []string {
	var p []string
	_, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to allocate root path.")
	}
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				p = append(p, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return p
}

func filterDirectory(path string, keyword string) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Unable to read file from the path: %v", err)
	}
	if fi.Mode().IsRegular() {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Unable to read file content: %v", err)
		}

		matched, err := regexp.Match(keyword, f)
		if err != nil {
			log.Fatalf("Unable to match the regexp pattern")
		}

		if matched {
			fmt.Println("File path: ", path)
		}
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

	scanner := bufio.NewScanner(f)

	line := 1

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
