package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	keywordsArgs := flag.String("keyword", "", "Keyword to search")
	absolutePath := flag.Bool("absolute", false, "Display absolute path")
	flag.Parse()

	if len(*keywordsArgs) == 0 {
		panic("--keyword is required")
	}

	var ignore map[string]struct{}

	if f, err := os.OpenFile(".folderignore", os.O_RDONLY, 0644); err == nil {
		ignore = readIgnoreFile(f)
		defer f.Close()
	}

	filePaths, err := scanDirectory(ignore)
	if err != nil {
		log.Fatalf("Failed to scan the project structure: %v", err)
	}

	showOutput(filePaths, *keywordsArgs, *absolutePath)

}

func readIgnoreFile(f *os.File) map[string]struct{} {
	ignore := map[string]struct{}{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		p := strings.TrimSpace(s.Text())
		d := strings.Split(p, "/")
		ignore[d[len(d)-1]] = struct{}{}
	}
	f.Close()
	return ignore
}

func scanDirectory(ignore map[string]struct{}) ([]string, error) {
	var paths []string

	err := filepath.Walk(".", func(path string, file fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, ok := ignore[file.Name()]

		if file.IsDir() && ok {
			return filepath.SkipDir
		}

		if file.Mode().IsRegular() && !ok {
			paths = append(paths, path)
		}

		return nil
	})

	return paths, err
}

func scanLine(writer io.Writer, path string, keyword string) []string {

	var loc []string

	f, _ := os.OpenFile(path, os.O_RDONLY, 0644)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 512*1024)
	line := 1

	scanner.Buffer(buf, 512*1024)

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read the text: %v", err)
	}

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			loc = append(loc, fmt.Sprintf("\tLine: %d", line))
		}
		line++
	}

	return loc
}

func showOutput(filePaths []string, keyword string, absPath bool) {
	for _, p := range filePaths {
		loc := scanLine(os.Stdout, p, keyword)
		if len(loc) > 0 {
			if absPath {
				root, err := os.Getwd()
				if err != nil {
					log.Fatalf("Failed to aollcate root path: %v", err)
				}
				p = path.Join(root, p)
			}
			fmt.Printf("%s: \n", p)
		}

		for _, line := range loc {

			fmt.Println(line)
		}
	}
}
