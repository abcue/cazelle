package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCueImports(t *testing.T) {
	// Setup temporary test directory
	testDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	tests := []struct {
		name     string
		files    map[string]string
		expected []string
	}{
		{
			name: "single_import",
			files: map[string]string{
				"test.cue": `import "k8s.io/api/core/v1"`,
			},
			expected: []string{"k8s.io/api/core/v1"},
		},
		{
			name: "aliased_import",
			files: map[string]string{
				"test.cue": `import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
			},
			expected: []string{"k8s.io/apimachinery/pkg/apis/meta/v1"},
		},
		{
			name: "factored_imports",
			files: map[string]string{
				"test.cue": `import (
					"list"
					foo "custom/pkg"
				)`,
			},
			expected: []string{"custom/pkg", "list"},
		},
		{
			name: "mixed_imports",
			files: map[string]string{
				"test1.cue": `import "strings"`,
				"test2.cue": `import (
					bar "another/pkg/v2"
					"math"
				)`,
			},
			expected: []string{"another/pkg/v2", "math", "strings"},
		},
		{
			name: "empty_file",
			files: map[string]string{
				"empty.cue": "",
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caseDir := filepath.Join(testDir, tt.name)
			os.MkdirAll(caseDir, 0755)
			os.Chdir(caseDir)
			defer os.Chdir(testDir)

			// Create test files with full paths
			for filename, content := range tt.files {
				fullPath := filepath.Join(caseDir, filename)
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			imports := findCueImports()

			// Verify results
			if len(imports) != len(tt.expected) {
				t.Fatalf("Expected %d imports, got %d:\n%v", len(tt.expected), len(imports), imports)
			}

			for i := range imports {
				if i >= len(tt.expected) {
					t.Fatalf("Unexpected extra import: %q", imports[i])
				}
				if imports[i] != tt.expected[i] {
					t.Errorf("Import mismatch at index %d:\nExpected: %q\nGot:      %q",
						i, tt.expected[i], imports[i])
				}
			}
			if len(imports) > len(tt.expected) {
				t.Errorf("Extra imports found: %v", imports[len(tt.expected):])
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	t.Run("valid_template", func(t *testing.T) {
		imports := []string{"pkg1", "pkg2"}
		templateContent := "Imports: {{imports}}"

		// Create temp template file
		path := filepath.Join(t.TempDir(), "template.txt")
		if err := os.WriteFile(path, []byte(templateContent), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := renderTemplate(imports, path)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "Imports: pkg1, pkg2"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("missing_template", func(t *testing.T) {
		_, err := renderTemplate([]string{}, "nonexistent.txt")
		if err == nil {
			t.Error("Expected error for missing template file")
		}
	})

	t.Run("empty_imports", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "template.txt")
		if err := os.WriteFile(path, []byte("{{imports}}"), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := renderTemplate([]string{}, path)
		if err != nil {
			t.Fatal(err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestMain(m *testing.M) {
	// Preserve original command-line arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Run tests
	code := m.Run()
	os.Exit(code)
}
