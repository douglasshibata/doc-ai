package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIProvider struct {
	APIKey  string
	BaseURL string
	Model   string
}

func NewOpenAIProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1" // Default
	}
	if model == "" {
		model = "gpt-4o"
	}
	return &OpenAIProvider{APIKey: apiKey, BaseURL: baseURL, Model: model}
}

func (p *OpenAIProvider) GenerateDocumentation(ctx context.Context, currentDoc, gitDiff string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", p.BaseURL)

	systemPrompt := "Você é um documentador técnico estrito. Atualize a documentação existente com base nas novas alterações de código fornecidas no Diff. Não remova as documentações existentes de funções que não foram alteradas. Adicione novas funções, atualize as modificações de parâmetros e ajuste o diagrama Mermaid.js se o fluxo de dados mudou. Retorne apenas o markdown atualizado."

	prompt := fmt.Sprintf("Documentação Atual:\n%s\n\nAlterações de Código (Diff):\n%s", currentDoc, gitDiff)

	body := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return result.Choices[0].Message.Content, nil
}
