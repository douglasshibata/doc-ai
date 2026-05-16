package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/doc-ai/internal/git"
	"github.com/user/doc-ai/internal/llm"
	"github.com/user/doc-ai/internal/orchestrator"
	"github.com/user/doc-ai/internal/parser"
)

var (
	providerName string
	apiKey       string
	baseURL      string
	modelName    string
	fromHash     string
	toHash       string
)

var rootCmd = &cobra.Command{
	Use:   "doc-ai",
	Short: "doc-ai is a tool to generate and update documentation automatically",
}

var incrementalCmd = &cobra.Command{
	Use:   "incremental [path]",
	Short: "Update documentation incrementally based on git diff",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		var provider llm.Provider
		switch providerName {
		case "openai", "ollama":
			provider = llm.NewOpenAIProvider(apiKey, baseURL, modelName)
		case "gemini":
			provider = llm.NewGeminiProvider(apiKey, modelName)
		default:
			return fmt.Errorf("unknown provider: %s", providerName)
		}

		if fromHash == "" || toHash == "" {
			var err error
			fromHash, toHash, err = git.GetLastTwoHashes(path)
			if err != nil {
				return err
			}
			fmt.Printf("Using commits %s -> %s\n", fromHash, toHash)
		}

		orch := orchestrator.NewOrchestrator(provider, path)
		return orch.RunIncremental(context.Background(), fromHash, toHash)
	},
}

var coverageCmd = &cobra.Command{
	Use:   "coverage [path]",
	Short: "Calculate documentation coverage",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		metrics, err := parser.CalculateCoverage(path)
		if err != nil {
			return err
		}

		fmt.Printf("Documentation Coverage: %.2f%%\n", metrics.Percentage)
		fmt.Printf("Total Exported Symbols: %d\n", metrics.TotalFunctions)
		fmt.Printf("Documented Symbols: %d\n", metrics.DocumentedFunctions)

		// Generate badge Markdown
		color := "red"
		if metrics.Percentage > 80 {
			color = "green"
		} else if metrics.Percentage > 50 {
			color = "yellow"
		}
		badge := fmt.Sprintf("![Documentation Coverage](https://img.shields.io/badge/Documentation_Coverage-%.2f%%25-%s)", metrics.Percentage, color)
		fmt.Printf("\nBadge Markdown:\n%s\n", badge)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&providerName, "provider", "gemini", "AI provider (gemini, openai, ollama)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", os.Getenv("GEMINI_API_KEY"), "API Key for the provider")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL for the provider (required for ollama)")
	rootCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name")

	incrementalCmd.Flags().StringVar(&fromHash, "from", "", "From commit hash")
	incrementalCmd.Flags().StringVar(&toHash, "to", "", "To commit hash")

	rootCmd.AddCommand(incrementalCmd)
	rootCmd.AddCommand(coverageCmd)
}
