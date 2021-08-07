package main

import (
	"bufio"
	"flag"
	"fmt"
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
	var keywords []string

	for _, k := range strings.Split(*keywordsArgs, ",") {
		keywords = append(keywords, strings.TrimSpace(k))
	}

	if f, err := os.OpenFile(".folderignore", os.O_RDONLY, 0644); err == nil {
		ignore = readIgnoreFile(f)
		defer f.Close()
	}

	filePaths, err := scanDirectory(ignore)
	if err != nil {
		log.Fatalf("Failed to scan the project structure: %v", err)
	}

	for _, k := range keywords {
		fmt.Println("--------------------------------------------")
		fmt.Printf("Keyword: %s\n\n", k)
		for _, p := range filePaths {
			isFound := false
			isDetectKeyword(*absolutePath, p, k, &isFound)
			if isFound {
				showLineContainsText(p, k)
			}
		}
	}

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
	var p []string

	err := filepath.Walk(".", traversal(ignore, &p))

	if err != nil {
		return nil, err
	}

	return p, nil
}

func traversal(ignore map[string]struct{}, p *[]string) func(path string, info os.FileInfo, err error) error {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, ok := ignore[file.Name()]

		if file.IsDir() && ok {
			return filepath.SkipDir
		}

		if file.Mode().IsRegular() && !ok {
			*p = append(*p, path)
		}

		return nil
	}
}

func isDetectKeyword(absolute bool, filePath string, keyword string, isFound *bool) {
	var p string

	if absolute {
		root, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to aollcate root path: %v", err)
		}
		p = path.Join(root, filePath)
	} else {
		p = filePath
	}

	f, _ := os.OpenFile(filePath, os.O_RDONLY, 0644)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 512*1024)

	scanner.Buffer(buf, 512*1024)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			*isFound = true
		}
	}

	if *isFound {
		fmt.Printf("%s\n", p)
	}
}

func showLineContainsText(path string, keyword string) {
	f, _ := os.OpenFile(path, os.O_RDONLY, 0644)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 512*1024)
	line := 1

	scanner.Buffer(buf, 512*1024)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			fmt.Printf("\tLine: %d\n", line)
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read the text: %v", err)
	}
}
