package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/doc-ai/internal/config"
	"github.com/user/doc-ai/internal/git"
	"github.com/user/doc-ai/internal/llm"
	"github.com/user/doc-ai/internal/parser"
)

type Orchestrator struct {
	LLM      llm.Provider
	RootPath string
}

func NewOrchestrator(provider llm.Provider, rootPath string) *Orchestrator {
	return &Orchestrator{
		LLM:      provider,
		RootPath: rootPath,
	}
}

func (o *Orchestrator) RunIncremental(ctx context.Context, fromHash, toHash string) error {
	state, err := config.LoadState(o.RootPath)
	if err != nil {
		return err
	}

	changedFiles, err := git.GetDiffBetweenCommits(o.RootPath, fromHash, toHash)
	if err != nil {
		return err
	}

	for _, cf := range changedFiles {
		if !strings.HasSuffix(cf.Path, ".go") || strings.HasSuffix(cf.Path, "_test.go") {
			continue
		}

		fullPath := filepath.Join(o.RootPath, cf.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			delete(state.Files, cf.Path)
			continue // File deleted
		}

		// Optional: use hash to skip if no real changes
		newHash, _ := config.CalculateHash(fullPath)
		if oldHash, ok := state.Files[cf.Path]; ok && oldHash == newHash {
			// No change since last documented (if state was updated)
			// But we trust Git diff more for incremental command
		}

		modifiedFuncs, err := parser.AnalyzeFile(fullPath, cf.Lines)
		if err != nil {
			fmt.Printf("Warning: failed to analyze file %s: %v\n", cf.Path, err)
			continue
		}

		if len(modifiedFuncs) == 0 {
			continue
		}

		// Load existing doc
		docDir := filepath.Join(o.RootPath, "docs")
		docPath := filepath.Join(docDir, strings.TrimSuffix(cf.Path, ".go")+".md")
		currentDoc := ""
		if data, err := os.ReadFile(docPath); err == nil {
			currentDoc = string(data)
		}

		// Prepare Diff for LLM
		var diffBuilder strings.Builder
		for _, f := range modifiedFuncs {
			diffBuilder.WriteString(fmt.Sprintf("Function/Type: %s\nCode:\n%s\n\n", f.Name, f.Code))
		}

		fmt.Printf("Updating documentation for %s...\n", cf.Path)
		updatedDoc, err := o.LLM.GenerateDocumentation(ctx, currentDoc, diffBuilder.String())
		if err != nil {
			return fmt.Errorf("LLM error for %s: %w", cf.Path, err)
		}

		// Save updated doc
		if err := os.MkdirAll(filepath.Dir(docPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(docPath, []byte(updatedDoc), 0644); err != nil {
			return err
		}

		state.Files[cf.Path] = newHash
	}

	return state.Save(o.RootPath)
}
