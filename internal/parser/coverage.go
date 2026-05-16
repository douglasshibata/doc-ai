package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type CoverageMetrics struct {
	TotalFunctions      int
	DocumentedFunctions int
	Percentage          float64
}

func CalculateCoverage(rootPath string) (CoverageMetrics, error) {
	metrics := CoverageMetrics{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (info.Name() == ".git" || info.Name() == ".docai" || info.Name() == "docs") {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files that don't parse
		}

		// Count exported functions and structs
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.IsExported() {
					metrics.TotalFunctions++
					if x.Doc != nil && len(strings.TrimSpace(x.Doc.Text())) > 0 {
						metrics.DocumentedFunctions++
					} else {
						// Check if there is a corresponding doc in docs/ folder
						if hasExternalDoc(rootPath, path, x.Name.Name) {
							metrics.DocumentedFunctions++
						}
					}
				}
			case *ast.TypeSpec:
				if x.Name.IsExported() {
					metrics.TotalFunctions++
					if x.Doc != nil && len(strings.TrimSpace(x.Doc.Text())) > 0 {
						metrics.DocumentedFunctions++
					} else {
						if hasExternalDoc(rootPath, path, x.Name.Name) {
							metrics.DocumentedFunctions++
						}
					}
				}
			}
			return true
		})

		return nil
	})

	if metrics.TotalFunctions > 0 {
		metrics.Percentage = float64(metrics.DocumentedFunctions) / float64(metrics.TotalFunctions) * 100
	}

	return metrics, err
}

func hasExternalDoc(rootPath, filePath, name string) bool {
	rel, _ := filepath.Rel(rootPath, filePath)
	docPath := filepath.Join(rootPath, "docs", strings.TrimSuffix(rel, ".go")+".md")

	data, err := os.ReadFile(docPath)
	if err != nil {
		return false
	}

	// Simple check: if the name appears in the markdown file
	return strings.Contains(string(data), name)
}
