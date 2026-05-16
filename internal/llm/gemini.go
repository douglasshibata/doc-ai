package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	APIKey string
	Model  string
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-1.5-flash"
	}
	return &GeminiProvider{APIKey: apiKey, Model: model}
}

func (p *GeminiProvider) GenerateDocumentation(ctx context.Context, currentDoc, gitDiff string) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.APIKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel(p.Model)
	model.SetTemperature(0.2)

	systemPrompt := "Você é um documentador técnico estrito. Atualize a documentação existente com base nas novas alterações de código fornecidas no Diff. Não remova as documentações existentes de funções que não foram alteradas. Adicione novas funções, atualize as modificações de parâmetros e ajuste o diagrama Mermaid.js se o fluxo de dados mudou. Retorne apenas o markdown atualizado."

	prompt := fmt.Sprintf("Documentação Atual:\n%s\n\nAlterações de Código (Diff):\n%s", currentDoc, gitDiff)

	// Note: Gemini API currently (in some versions) doesn't have a dedicated "System Instruction"
	// in the same way as OpenAI, but we can prepend it or use the dedicated field if available.
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from Gemini")
	}

	var content strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content.WriteString(string(text))
		}
	}

	return content.String(), nil
}
