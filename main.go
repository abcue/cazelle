// Co-authored by https://www.perplexity.ai/search/get-all-imported-cuelang-packa-TnTtLs06Q5CM_bGpcfHMqA
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func findCueImports() []string {
	importSet := make(map[string]struct{})
	importPattern := regexp.MustCompile(`import\s*\(\s*([^)]+)\s*\)|import\s+"([^"]+)"`)
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal("Directory read error:", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cue") {
			content, err := os.ReadFile(file.Name())
			if err != nil {
				log.Printf("Error reading %s: %v", file.Name(), err)
				continue
			}
			matches := importPattern.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				groups := strings.ReplaceAll(strings.TrimSpace(match[1]+match[2]), "\n", " ")
				for _, imp := range strings.Split(groups, " ") {
					imp = strings.Trim(imp, `"`)
					if imp != "" {
						importSet[imp] = struct{}{}
					}
				}
			}
		}
	}
	imports := make([]string, 0, len(importSet))
	for imp := range importSet {
		imports = append(imports, imp)
	}
	sort.Strings(imports)
	return imports
}
func main() {
	fmt.Println("Imported packages:", findCueImports())
}
