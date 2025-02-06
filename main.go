package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func findCueImports() []string {
	importSet := make(map[string]struct{})
	importPattern := regexp.MustCompile(`(?m)(?:import\s*\(\s*([^)]+)\s*\))|(?:import\s+(?:[\w\.]+\s+)?\"([^\"]+)\")`)

	files, err := os.ReadDir(".")
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
				if match[1] != "" { // Factored imports
					for _, imp := range strings.Split(match[1], "\n") {
						imp = strings.TrimSpace(imp)
						if imp == "" {
							continue
						}
						parts := strings.Fields(imp)
						if len(parts) == 0 {
							continue
						}
						path := strings.Trim(parts[len(parts)-1], `"`)
						if path != "" {
							importSet[path] = struct{}{}
						}
					}
				} else if match[2] != "" { // Single import (aliased or direct)
					path := strings.Trim(match[2], `"`)
					if path != "" {
						importSet[path] = struct{}{}
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

func renderTemplate(imports []string, templatePath string) (string, error) {
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %w", err)
	}

	importsList := strings.Join(imports, ", ")
	rendered := strings.ReplaceAll(string(templateContent), "{{imports}}", importsList)
	return rendered, nil
}

func main() {
	templatePath := flag.String("template", "", "Path to template file")
	flag.Parse()

	imports := findCueImports()

	if *templatePath != "" {
		result, err := renderTemplate(imports, *templatePath)
		if err != nil {
			log.Fatal("Rendering error:", err)
		}
		fmt.Print(result)
	} else {
		fmt.Println("Detected imports:")
		for _, imp := range imports {
			fmt.Println("-", imp)
		}
	}
}
