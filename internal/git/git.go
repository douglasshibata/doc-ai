package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type ChangedFile struct {
	Path  string
	Lines []int
}

// GetDiffBetweenCommits returns a list of files and the lines that were changed between two commits.
func GetDiffBetweenCommits(repoPath string, fromHash, toHash string) ([]ChangedFile, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %w", err)
	}

	from, err := repo.ResolveRevision(plumbing.Revision(fromHash))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve fromHash: %w", err)
	}

	to, err := repo.ResolveRevision(plumbing.Revision(toHash))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve toHash: %w", err)
	}

	fromCommit, err := repo.CommitObject(*from)
	if err != nil {
		return nil, err
	}

	toCommit, err := repo.CommitObject(*to)
	if err != nil {
		return nil, err
	}

	patch, err := fromCommit.Patch(toCommit)
	if err != nil {
		return nil, err
	}

	var changedFiles []ChangedFile
	for _, fp := range patch.FilePatches() {
		_, toFile := fp.Files()
		if toFile == nil {
			continue // Deleted file
		}

		path := toFile.Path()
		changedFile := ChangedFile{Path: path}

		// Track line numbers in the NEW version of the file
		currentLine := 1
		for _, chunk := range fp.Chunks() {
			content := chunk.Content()
			lines := strings.Split(content, "\n")
			// Last split might be empty if content ends in \n
			numLines := len(lines)
			if content != "" && !strings.HasSuffix(content, "\n") {
				// numLines is correct
			} else if content != "" {
				numLines--
			} else {
				numLines = 0
			}

			switch chunk.Type() {
			case 0: // Equal (plumbing/object/patch.go: const Equal DeltaOperation = 0)
				currentLine += numLines
			case 1: // Add (plumbing/object/patch.go: const Add DeltaOperation = 1)
				for i := 0; i < numLines; i++ {
					changedFile.Lines = append(changedFile.Lines, currentLine+i)
				}
				currentLine += numLines
			case 2: // Delete (plumbing/object/patch.go: const Delete DeltaOperation = 2)
				// Deletions mark the current line as changed in context
				changedFile.Lines = append(changedFile.Lines, currentLine)
			}
		}
		changedFiles = append(changedFiles, changedFile)
	}

	return changedFiles, nil
}

// GetLastTwoHashes returns the HEAD and HEAD~1 hashes.
func GetLastTwoHashes(repoPath string) (string, string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", "", err
	}

	currentHash := head.Hash().String()

	// Get previous commit
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", "", err
	}

	if commit.NumParents() == 0 {
		return "", currentHash, nil // First commit
	}

	parent, err := commit.Parent(0)
	if err != nil {
		return "", currentHash, nil
	}

	return parent.Hash.String(), currentHash, nil
}
