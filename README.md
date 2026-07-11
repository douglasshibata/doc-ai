# doc-ai: Incremental Documentation CLI

`doc-ai` is a Go-based tool designed to evolve your project's documentation at each commit. It uses Git diffs and AST analysis to identify code changes and uses AI (Gemini/OpenAI) to update only the relevant parts of your documentation.

## 🚀 Features

- **Incremental Updates:** Only sends modified functions/structs to the LLM.
- **AST-Aware:** Pinpoints changes using Go's Abstract Syntax Tree.
- **Multi-Provider:** Supports Google Gemini, OpenAI, and Ollama.
- **Coverage Metrics:** Tracks documentation progress and generates status badges.
- **CI/CD Ready:** Templates for GitHub Actions, Bitbucket Pipelines, and Jenkins.

## 🛠️ Installation

```bash
go install github.com/user/doc-ai@latest
```

## 📖 Usage

### Initialize
The tool tracks state in a `.docai/` folder.
```bash
doc-ai incremental .
```

### Check Coverage
Generate a documentation coverage report and a status badge.
```bash
doc-ai coverage .
```

### Configuration
Set your API keys via environment variables:
```bash
export GEMINI_API_KEY="your-key"
# OR
export OPENAI_API_KEY="your-key"
```

## 🤖 CI/CD Integration

### Bitbucket Pipelines
Use the provided `bitbucket-pipelines.yml`:
- Set `GEMINI_API_KEY` in Repository Variables.
- Ensure the Pipeline has write access to the repository.

### Jenkins
Use the provided `Jenkinsfile`:
- Configure `gemini-api-key` as a Secret Text credential.
- Use `[skip ci]` in commit messages to avoid loops.

## 🏗️ Architecture

```text
[Commit] -> [Git Diff] -> [AST Parser] -> [LLM Provider] -> [Markdown Update]
```

---
*Powered by AI to keep your code and docs in sync.*

---

## 🛠️ Project Deep Dive & Technical Guide

### 🎯 Main Purpose
Automated Documentation CLI generator for multi-language projects.

### 🧰 Tecnologias Utilizadas (Tech Stack)
- **Primary Technologies:** Go 1.24.3
- **Deployment Strategy:** Container-ready with standard Docker/Docker Compose integrations where applicable.

### 🚀 Como Implantar e Executar (How to Deploy & Run)
#### Go CLI / Service
1. Build the module: `go build`.
2. Execute: `./doc-ai` (or project executable name).

---
