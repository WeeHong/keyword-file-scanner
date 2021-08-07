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
	keywords := flag.String("keyword", "", "Keyword to search")
	absolutePath := flag.Bool("absolute", false, "Display absolute path")
	flag.Parse()

	if len(*keywords) == 0 {
		panic("--keyword is required")
	}

	ignore := readIgnoreFile()

	kys := strings.TrimSpace(*keywords)
	for _, k := range strings.Split(kys, ",") {

		filePaths := scanDirectory(ignore, k)
		fmt.Printf("Keyword \"%s\" search resullt:\n", k)

		for _, file := range filePaths {
			if *absolutePath {
				root, err := os.Getwd()
				if err != nil {
					log.Fatalf("Unable to allocate root path: %v", err)
				}
				fmt.Println(path.Join(root, file))
			} else {
				fmt.Println(file)
			}
		}
	}
}

func readIgnoreFile() map[string]struct{} {
	ignore := map[string]struct{}{}
	if f, err := os.OpenFile(".folderignore", os.O_RDONLY, 0644); err == nil {
		s := bufio.NewScanner(f)
		for s.Scan() {
			p := strings.TrimSpace(s.Text())
			d := strings.Split(p, "/")
			ignore[d[len(d)-1]] = struct{}{}
		}
		f.Close()
	}
	return ignore
}

func scanDirectory(ignore map[string]struct{}, keyword string) []string {
	var p []string

	err := filepath.Walk(".", traversal(ignore, &p, keyword))

	if err != nil {
		log.Println(err)
	}

	return p
}

func traversal(ignore map[string]struct{}, p *[]string, keyword string) func(path string, info os.FileInfo, err error) error {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, ok := ignore[file.Name()]

		if file.IsDir() && ok {
			return filepath.SkipDir
		}

		if file.Mode().IsRegular() && !ok {
			showLineContainsText(path, keyword, p)
		}

		return nil
	}
}

func showLineContainsText(path string, keyword string, p *[]string) {
	f, _ := os.OpenFile(path, os.O_RDONLY, 0644)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 512*1024)
	scanner.Buffer(buf, 512*1024)

	line := 1

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			s := fmt.Sprintf("Found at: %s, Line: %d", path, line)
			*p = append(*p, s)
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read the text: %v", err)
	}
}
