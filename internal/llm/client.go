package llm

import "context"

// Provider defines the contract for any AI client in the CLI
type Provider interface {
	GenerateDocumentation(ctx context.Context, currentDoc, gitDiff string) (string, error)
}
