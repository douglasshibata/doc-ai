package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// FuncDiff stores the cleaned modified structures
type FuncDiff struct {
	Name string
	Doc  string
	Code string
}

// AnalyzeFile analyzes the file focusing only on the modified structs and functions
func AnalyzeFile(filePath string, changedLines []int) ([]FuncDiff, error) {
	fset := token.NewFileSet()

	// Read file from disk
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the AST
	node, err := parser.ParseFile(fset, filePath, fileContent, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var modifiedFuncs []FuncDiff

	// Inspect the file's tree
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Get start and end lines of the function
			startLine := fset.Position(x.Pos()).Line
			endLine := fset.Position(x.End()).Line

			// Check if any of the lines changed by Git are within this function's scope
			if isIntersecting(startLine, endLine, changedLines) {
				// Extract clean raw function code
				funcCode := string(fileContent[x.Pos()-1 : x.End()])

				docText := ""
				if x.Doc != nil {
					docText = x.Doc.Text()
				}

				modifiedFuncs = append(modifiedFuncs, FuncDiff{
					Name: x.Name.Name,
					Doc:  docText,
					Code: funcCode,
				})
			}
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				startLine := fset.Position(x.Pos()).Line
				endLine := fset.Position(x.End()).Line

				if isIntersecting(startLine, endLine, changedLines) {
					for _, spec := range x.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}

						typeName := typeSpec.Name.Name
						code := string(fileContent[x.Pos()-1 : x.End()])
						doc := ""
						if x.Doc != nil {
							doc = x.Doc.Text()
						} else if typeSpec.Doc != nil {
							doc = typeSpec.Doc.Text()
						}

						modifiedFuncs = append(modifiedFuncs, FuncDiff{
							Name: typeName,
							Doc:  doc,
							Code: code,
						})
					}
				}
			}
		}
		return true
	})

	return modifiedFuncs, nil
}

func isIntersecting(start, end int, changedLines []int) bool {
	for _, line := range changedLines {
		if line >= start && line <= end {
			return true
		}
	}
	return false
}
